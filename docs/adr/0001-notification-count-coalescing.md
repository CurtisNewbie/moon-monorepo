# Coalesce unread-count pushes per user per second

Status: accepted

A bulk "Mark All Opened" updates N notification rows in one statement, and MySQL's binlog reports each updated row separately. The notification change stream forwards each row as its own message, so one user action fans out to N cache evictions, N Redis publishes, N unread-count recomputations and N WebSocket frames — all carrying the same final count. We coalesce the publisher side: the replica consuming the change stream claims a per-user Redis key with `SET NX EX <window>` and schedules one trailing publish a window later; every other event in that window is dropped. The announcement carries only the user id, the channel's existing payload, because every subscriber recomputes the count from the database — so no value needs to travel and no subscriber changes. The window defaults to one second and is configurable.

## Considered Options

- **In-process coalescing at the stream consumer.** Collapses N to one only while a single replica consumes the stream. The stream is a shared durable queue with competing consumers, so with R replicas a user's events are spread across replicas and each collapses a fragment — up to R publishes per window. A Redis-owned window holds at any R.
- **Coalescing at the Redis subscriber, before recomputation.** Leaves all N publishes on the wire, which is the bulk of the waste.
- **Recognising the bulk statement and collapsing its rows as a unit.** Not possible: the streamed event carries no statement identity — no transaction id, log position or affected-row count — only per-row columns and a whole-second timestamp.

## Consequences

- Every count change is up to one window late.
- A replica that dies inside its window loses that window's announcement. The key expires on its own, so the next change opens a new window, and clients recover on reconnect, on the long-poll refresh, or on the next change. No durable queue is introduced.
- Dropping intermediate counts costs nothing, because the count is a single integer recomputed from the database and only its latest value is meaningful.
- A burst longer than one window keeps announcing one change per window; that is throttling, not a leak.
