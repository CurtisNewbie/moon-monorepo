# Moon Platform

The platform in this monorepo: a file manager (`vfm`), user notifications (`user-vault`), and the Angular client that fronts them.

## Language

### Platform Notifications

**Login Session**:
The period during which the client holds a confirmed identity for a User. Unread Notification Count is only ever pushed to a client inside one.
_Avoid_: authenticated state, auth session, logged-in flag

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

### Layout

**Compact Layout**:
The layout a client uses when its viewport is 768px wide or narrower. It changes presentation, the column list a table shows, and the shape of an interaction; it never removes a capability.
_Avoid_: mobile layout, phone view, small screen

**Wide Layout**:
The layout a client uses when its viewport is 769px wide or wider. A phone held sideways reaches it, and so does a narrow desktop window that grew past the boundary.
_Avoid_: desktop layout, full view, large screen

### File Browsing

**Resume Cursor**:
The file a user last read in a directory, remembered per user, per directory and per sort order, so that reopening the directory under the same sort returns them there.
_Avoid_: Last page, bookmark, reading position, progress

**Order-By Preference**:
The sort order a user last chose for a directory. A directory opens in its Order-By Preference.
_Avoid_: Sort setting, default sort, sort cache

**Size Sort**:
A listing ordered by file size, offered at all times and files only, so that a directory row, which has no size, never sits in it. Never stored as an Order-By Preference and never resumed. At root it spans every directory.
_Avoid_: sort by bytes, size order, biggest first

**Search Term**:
One fragment a Keyword Search requires a name to contain. What the user types is split into Search Terms at every separator, at a letter-to-digit or lower-to-upper change, and, for a run of Latin letters, into dictionary words, so `my holiday`, `my-holiday` and `myholiday` all yield the same two terms. Every term must be present, in any order.
_Avoid_: token, keyword word, n-gram

**Keyword Search**:
A listing narrowed to files whose name contains every Search Term. Inside a directory it covers that directory alone; at root it covers every directory.
_Avoid_: fulltext search, fuzzy search, filter
