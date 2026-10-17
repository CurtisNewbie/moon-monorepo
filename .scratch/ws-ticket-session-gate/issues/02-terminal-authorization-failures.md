# 02: Authorization failures terminal for the session

**What to build:** Inside a Login Session, a ticket request refused with 401 or 403 ends notification subscription for that session: no retry is scheduled and no further request is made until the session changes. Transport failures — network errors and 5xx — still retry on the existing exponential backoff from one second to a thirty-second ceiling. An in-session socket close caused by transport failure reconnects under the same gate and the same terminal rule; a session-end close never reconnects.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] A 401 or 403 on the ticket request schedules no retry and produces no further ticket request for that session
- [ ] A network error or 5xx on the ticket request retries after the backoff delay
- [ ] A transport-failure socket close inside the session reconnects, still subject to the gate
- [ ] Re-login after a terminal failure starts a fresh connection attempt
- [ ] Tests cover the above at the service's public seam and pass under the app's Karma/ChromeHeadless suite

## Comments
