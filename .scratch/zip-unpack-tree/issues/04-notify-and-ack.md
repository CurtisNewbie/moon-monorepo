# 04: Tell the user what happened, and never retry a failed unpack

**What to build:** An unpack that drops entries, is refused, or stops early leaves the user with one notification saying what happened. A refusal also fails the action itself, so it is visible on the click as well as in the inbox.

A write failure mid-archive stops the remaining work, is logged at ERROR level, and is acknowledged rather than returned, so the message is never redelivered. The partial tree stays. Sending a notification that fails is logged and never changes the outcome.

**Blocked by:** 01, 02

**Status:** ready-for-agent

- [ ] A completed unpack that dropped entries sends one notification naming how many were dropped and which paths.
- [ ] A refused archive sends one notification naming the reason, and the action returns the refusal error to the caller.
- [ ] An unpack that stops early sends one notification with how many entries were not created, and its handler returns no error.
- [ ] A run that both dropped entries and stopped early sends one notification, not two.
- [ ] A run that dropped nothing, was not refused, and did not stop early sends no notification.
- [ ] The notification goes through user-vault with a title within 255 characters and a message within 1000, and a message longer than the limit is shortened rather than dropped.
- [ ] A failed notification call is logged and never fails the unpack.
- [ ] No retry header is added, and the handler acknowledges on write failure.
