# Notification WebSocket Only Inside a Login Session

Status: ready-for-agent

## Problem Statement

A logged-out visitor to the moon client is treated as a notification subscriber.

The Angular app constructs the notification WebSocket service during bootstrap, before anyone has logged in, and that service asks user-vault for a websocket ticket immediately. With no token the request is rejected, and each rejection is handled twice: the websocket service schedules a retry with exponential backoff (capped at 30 seconds, forever), and the shared HTTP error interceptor treats the 401 as "the user's session died" — it shows a "Please login first" snackbar and calls logout, which navigates to the login page again.

The visible result is a login page that keeps announcing "Please login first" while the user is doing nothing, an endless stream of `ws-ticket` requests in the network panel, and a logout being invoked repeatedly against a session that never existed.

## Solution

Notifications subscribe only inside a Login Session.

A Login Session is established when the client's identity is confirmed by user-vault, not merely when a token happens to be in local storage. While no Login Session exists, the client makes no websocket ticket request and holds no socket. When a Login Session starts, the client connects once, fresh, with no attempt to replay whatever was skipped. When the Login Session ends, the socket is closed, pending retries are cancelled, the backoff state is discarded, and the badge returns to zero.

Being subscribed is a privilege of the session, so authorization failures are terminal for that session: a ticket request rejected with 401 or 403 stops the client from asking again until the session changes. Transport failures — network errors, 5xx — are still retried with the existing backoff.

## User Stories

1. As a logged-out visitor, I want the client to make no notification websocket request, so that my login page is not polluted with requests for something I cannot have.
2. As a logged-out visitor, I want no "Please login first" snackbar while I am reading the login page, so that I am not told to do the thing I am already doing.
3. As a logged-out visitor, I want opening the app to not trigger a logout, so that the app does not bounce me around before I have a session.
4. As a logged-out visitor who registers, I want the same quiet behavior, so that registration is not disturbed by notification traffic.
5. As a returning user with a valid token, I want the client to connect once my identity is confirmed, so that my unread Notification Count starts arriving without me refreshing.
6. As a returning user, I want exactly one websocket connection, so that the count is not delivered twice or fought over by duplicate sockets.
7. As a returning user, I want the connection to carry a fresh ticket, so that a ticket is never reused across sessions.
8. As a user who logs out, I want the socket closed and its retries cancelled, so that no connection outlives my session.
9. As a user who logs out, I want the unread badge cleared, so that the previous session's count is not shown to the next person at the machine.
10. As a user who logs out while a retry is pending, I want that retry discarded, so that the socket does not quietly come back after logout.
11. As a user whose token expired mid-session, I want nothing further requested after the rejection, so that an expired session does not generate an endless stream of failures.
12. As a user lacking notification permission, I want the client to stop asking after it is refused, so that a permanent refusal does not repeat every 30 seconds.
13. As a logged-in user on a flaky network, I want ticket retries to still back off and recover, so that a transient failure does not cost me notifications for the rest of the session.
14. As a logged-in user who reconnects after a dropped socket, I want the same session-level rules to apply, so that reconnecting never bypasses the session gate.
15. As a logged-in user, I want logging back in after a logout to re-establish the connection, so that the session gate is not a one-way door.

## Implementation Decisions

- Gate: a Login Session exists exactly when the client holds a confirmed identity — the user-info stream used elsewhere in the app emitting a non-null user. A token in local storage is not sufficient, because it can be stale or expired.
- Gate location: inside the notification websocket service itself, in one place, so its consumers are unchanged and there is a single choke point. Consumers must not need to know about the gate.
- Session start: connect once, fresh, with attempt counting reset. Nothing is queued or replayed — the unread count is current state, not an event log.
- Session end: close the socket, cancel the pending retry, reset the backoff state, and publish a count of zero.
- Stale retries: the connect path re-checks the gate, so a retry scheduled before logout cannot open a socket after logout.
- Terminal failures: a ticket request answered with 401 or 403 ends notification subscription for the current Login Session; no retry is scheduled and no further request is made until the session changes. 401 already triggers the app-wide logout; the terminal state is the guarantee that nothing retries behind it.
- Retryable failures: network errors and 5xx keep the existing exponential backoff from one second to a ceiling of thirty seconds.
- Socket lifecycle: an in-session close caused by transport failure still reconnects with the same backoff, subject to the same gate and terminal-failure rules; a session-end close is intentional and must not reconnect.
- No initial count hydration: the badge stays zero until the first Count-Changed Push of the session. Fetching the count over REST on connect is a separate concern.
- Untouched: the HTTP error interceptor, the empty-token callback, and their register-page special case. With the gate in place they no longer receive notification-driven 401s.
- No backend change: `ws-ticket` keeps requiring the notification query resource.
- Domain vocabulary: `GLOSSARY.md` gains Login Session, since the boundary being enforced is a domain boundary, not an implementation one.

## Testing Decisions

- Good test: assert externally observable behavior of the notification websocket service — which HTTP calls it makes, whether a socket is opened or closed, what closes it, and what the count stream emits. No assertions on private fields, timer identities, or internal method names.
- Seam: the service's existing public surface (its injected HTTP client and the user-info stream, its count observable, its destroy hook). No new seam is introduced, and no other module is modified to make testing possible.
- Test double for the browser's websocket: a fake installed on the window, so the test can observe constructed URLs and drive open/message/close events.
- Cases to cover: no ticket request before a confirmed user; exactly one ticket request and one socket after a confirmed user; the socket URL carries the returned ticket; user becoming null closes the socket, zeroes the count, and prevents reconnection; a retry scheduled before logout does not fire a request after logout; 401 and 403 schedule no retry; a 5xx retries after the backoff delay; a pushed count reaches the count observable; logging back in after logout re-requests a ticket.
- Prior art: the existing browse-history spec in the same app — `HttpClientTestingModule` plus a fake snackbar — is the pattern to follow for HTTP-facing services.
- Runner: the app's Karma/Jasmine suite in ChromeHeadless (`npm test`). The spec file is added by hand; the project's generator skips spec creation.

## Out of Scope

- Any browser or end-to-end verification of the real login page.
- Hydrating the unread count over REST when a session starts.
- Server-side enforcement of `ws-ticket` authorization.
- The v2 long-polling count endpoint, and any transport other than the websocket.
- Nav menu visibility, badge styling, and which routes require a session.
- Whether an expired token should be refreshed instead of ending the session.

## Further Notes

- The notification websocket service is root-scoped and never destroyed in practice, so cleanup that today exists only in its destroy hook must also run on session end.
- The project's Angular generator is configured to skip test generation; the new spec is written directly.
