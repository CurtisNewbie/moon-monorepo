# 01: Notification subscription gated by Login Session

**What to build:** A logged-out visitor's client makes no notification websocket ticket request and holds no socket. Once the client's identity is confirmed, it fetches one ticket and opens one socket; the unread count pushed over that socket reaches the badge stream. When the Login Session ends, the socket is closed, any pending retry is discarded, the backoff state is reset, and the count stream emits zero. A retry scheduled before the session ended must not open a socket after it ended.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] No ticket request is made before a confirmed identity, including on the login and register pages
- [x] A confirmed identity produces exactly one ticket request and one websocket connection, carrying that ticket
- [x] The identity stream becoming null closes the socket and emits a count of zero
- [x] A retry pending at session end fires no ticket request afterwards
- [x] Logging back in after logout establishes the connection again, once
- [x] Tests cover the above at the service's public seam and pass under the app's Karma/ChromeHeadless suite

## Comments
