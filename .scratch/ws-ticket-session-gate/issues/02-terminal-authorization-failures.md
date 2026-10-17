# 02: Authorization failures terminal for the session

**What to build:** Inside a Login Session, a ticket request refused with 401 or 403 ends notification subscription for that session: no retry is scheduled and no further request is made until the session changes. Transport failures — network errors and 5xx — still retry on the existing exponential backoff from one second to a thirty-second ceiling. An in-session socket close caused by transport failure reconnects under the same gate and the same terminal rule; a session-end close never reconnects.

**Blocked by:** 01

**Status:** done

- [x] A 401 or 403 on the ticket request schedules no retry and produces no further ticket request for that session
- [x] A network error or 5xx on the ticket request retries after the backoff delay
- [x] A transport-failure socket close inside the session reconnects, still subject to the gate
- [x] Re-login after a terminal failure starts a fresh connection attempt
- [x] Tests cover the above at the service's public seam and pass under the app's Karma/ChromeHeadless suite

## Comments
