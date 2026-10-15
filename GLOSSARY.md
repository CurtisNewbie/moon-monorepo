# Platform Notifications

The notification postbox in the `user-vault` service: per-user messages and their unread count, delivered to clients over WebSocket (with a long-poll fallback).

## Language

**Notification**:
A message addressed to exactly one recipient user, which is either unopened or opened.
_Avoid_: Postbox message, inbox item, alert

**Opened**:
The status of a Notification the recipient has seen. A Notification starts unopened and becomes opened one at a time, or in bulk.
_Avoid_: Read, seen, viewed, consumed

**Unread Notification Count**:
How many of a user's Notifications are still unopened. This single integer, per user, is the only thing pushed to a client.
_Avoid_: Badge count, notification total, inbox count

**Count-Changed Push**:
A delivery to one user's client carrying only that user's current Unread Notification Count.
_Avoid_: Count event, notification push, count message

**Coalescing Window**:
The period during which a user's repeated count changes are folded into a single Count-Changed Push. It opens on the first change for that user and defaults to one second.
_Avoid_: Throttle interval, debounce, batch
