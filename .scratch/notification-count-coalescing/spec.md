Status: ready-for-agent

# Notification count coalescing

## Problem Statement

When a user clicks "Mark All Opened", their unopened notifications are marked opened in a single statement covering many rows. MySQL's binlog reports each updated row separately and the notification change stream forwards each row as its own message, so that one action produces one cache eviction, one Redis publish, one unread-count recomputation and one WebSocket frame **per notification** — all carrying the same final count. Opening 1000 notifications at once means 1000× that work, caused entirely by the one user who clicked, for a count that never differed across those 1000 announcements.

A single-minded burst of work like this is self-inflicted load: it multiplies database queries, Redis traffic and socket writes by the number of notifications the user opened, for a change the user experiences as one event.

## Solution

Coalesce the publisher side so that at most one unread-count change is announced per user per coalescing window, defaulting to one second. The first change for a user opens a window; every other change in that window is dropped; exactly one announcement fires at the end of the window carrying the user's current count. The recipient's unread badge still converges to the correct value, a single change still reaches the client within the window, and a burst of any size collapses to one announcement per window.

## User Stories

1. As a notification recipient, I want my unread badge to show the correct count after I mark all notifications opened, so that I am not shown a stale number.
2. As a notification recipient, I want the badge to update within about a second of my action, so that the interface still feels responsive.
3. As a notification recipient, I want the badge to be correct when I load the page, even if my last change fell inside a window, so that I never start from a wrong count.
4. As a notification recipient, I want a single change (one open, or one new notification) to reach me within one window, so that everyday use is not noticeably delayed.
5. As a notification recipient, I want a new notification that arrives right after a burst not to be swallowed, so that an important message is never silently missed.
6. As a notification recipient, I want the badge to self-correct after a missed update, so that a rare loss is never permanent.
7. As a notification recipient using the long-poll fallback, I want the same coalesced behaviour as the WebSocket path, so that both transports agree.
8. As a platform operator, I want marking all of a user's notifications opened to cost one publish, not one per notification, so that one user's action cannot spike the service.
9. As a platform operator, I want it to cost one count recomputation, not one per notification, so that the database is not hammered.
10. As a platform operator, I want it to cost one frame per connection, not one per notification, so that socket writes stay proportional to connections, not to notifications.
11. As a platform operator, I want the coalescing window tunable by configuration, so that I can trade latency for load without a code change or a deploy.
12. As a platform operator, I want the behaviour to stay correct when the service runs as more than one replica, so that scaling out for capacity does not multiply announcements.
13. As an on-call engineer, I want a replica dying mid-window to lose at most one announcement, so that the failure mode is bounded and self-healing.
14. As an on-call engineer, I want a crashed replica never to suppress a user indefinitely, so that a stale window cannot wedge one user's badge forever.
15. As a developer, I want the announcement payload unchanged, so that the WebSocket and long-poll contracts do not move.
16. As a developer, I want the coalescing logic testable without a live Redis, so that the window rules are pinned by fast unit tests.
17. As a developer, I want the Redis gate covered by an integration test in the repo's existing style, so that the real `NX`/expiry behaviour is proven, not assumed.
18. As a developer, I want an end-to-end check that seeds a large unopened set, marks all, and observes exactly one announcement, so that the whole real path is verified.
19. As a developer, I want no schema change and no new endpoint, so that the change is contained to the publisher side.
20. As a developer, I want the join between the coalescer and the change stream to be a single named seam, so that future changes to the stream do not silently bypass coalescing.

## Implementation Decisions

- **Insertion point.** The notification change-stream listener in the postbox module — the handler that announces a count change for a user. It is the single publisher-side point; downstream subscribers are untouched in contract.
- **Window guard.** A per-user Redis key claimed with `SET NX EX <window>`. The replica that claims it owns the window; a replica that fails to claim drops its event. The guard lives in Redis rather than process memory because the change stream is a shared durable queue with competing consumers: with more than one replica a user's events land on different replicas, and only a shared gate sees the whole burst. The key's time-to-live equals the window, so expiry needs no separate cleanup.
- **Trailing flush.** The window owner schedules exactly one announcement `<window>` after it claims. The flush evicts the user's count cache first, then publishes the user id to the existing count-changed channel — eviction before publish so subscribers' recomputation reads current data. If that eviction fails the flush publishes nothing: announcing a change whose cache is still stale would make subscribers recompute and push a wrong count that survives until the cache expires.
- **No value in the announcement.** Subscribers recompute the count from the database, so the announcement carries only the user id — the channel's existing payload. This is what makes a Redis-owned window work: the owner never needs to know the latest value, only that a window is open.
- **One subscriber-side guard.** When a user has no live WebSocket connection, skip the count recomputation and the push. Long-poll subscribers are still woken and recompute for themselves.
- **Configuration.** The window is `user-vault.notification-count-coalesce-window`, a duration with a default of one second, declared as a `user-vault.*` config constant with the repo's property annotation and a `miso.SetDefProp` default so the service runs unchanged when the property is absent. Coalescing is always on; there is no off switch.
- **Contracts.** No new endpoint, no schema change, no change to the count-changed channel payload, no change to the WebSocket message shape, no change to the long-poll response shape.
- **Sustained bursts.** A burst spanning k windows produces k announcements, one per window — throttling, not a leak.

## Testing Decisions

- **What makes a good test here.** Tests assert observable behaviour — "a burst of N changes for one user yields exactly one announcement" — never the internal timer, map or keys.
- **Primary seam (new, unit, no infrastructure).** A count coalescer constructed with three injectable effects: a window gate (claim-if-absent with expiry), a flush action, and a clock/scheduler. Tests drive it with a synthetic burst and advance time. This is a new seam; no existing seam isolates the publish, so it is proposed at the highest usable point — the coalescer's input, one step below the change-stream listener. Cases it must pin: N changes in one window produce one flush; a burst spanning k windows produces k flushes; the gate is consulted for every change but claims at most once per window; the flush evicts before it publishes and publishes nothing when the eviction fails; a non-positive window is clamped so it can never become a zero (no-expiry) gate TTL.
- **Integration seam (repo convention, live Redis).** The Redis-backed gate exercised against a live Redis, loaded the way the service's existing integration tests load their dependencies. Verifies the key name, the `NX` claim semantics and the TTL expiry for real, since that behaviour is what the unit seam fakes. Prior art: the `preTest`/`preUserTest` setup in the user-vault service's existing tests.
- **Wired integration (live Redis).** Drive the real listener entry point with a burst of synthetic change events for one user, subscribe to the count-changed channel, and assert exactly one announcement arrives once the window has elapsed; a single event yields exactly one announcement. This covers the real gate, the real flush and the real publish. The binlog-to-event-bus path itself is not covered — it needs the event-pump process running — so the bulk update is not exercised end to end.
- **Regression.** A single open and a newly created notification each still produce an announcement within one window — coalescing must not suppress isolated changes. Existing pure unit tests in the service remain the model for the fast, infra-free cases.

## Out of Scope

- The bulk-mark rule itself. "Mark All Opened" marks rows up to the newest listed notification's id, so unopened rows outside the current page are left unopened. Deliberately left as-is.
- Suppressing no-op change events. Re-opening an already-opened notification still emits a change event; filtering it at the stream is not part of this work.
- The client. No frontend change; the badge keeps consuming the same frame.
- Observability for coalescing. No metric or log line reports dropped changes.

## Further Notes

- Why a Redis guard and not in-memory: the change stream is a single shared durable queue with competing consumers, so with more than one replica a user's events are distributed across replicas and a per-process coalescer sees only a fragment.
- Why the announcement carries no value: across replicas the window owner cannot know the latest count, and it does not need to — every subscriber recomputes it.
- Why the bulk statement cannot be recognised directly: the streamed event carries no statement identity (no transaction id, log position or affected-row count), only per-row columns and a whole-second timestamp.
- Rationale and rejected alternatives are recorded in the ADR on notification count coalescing.
