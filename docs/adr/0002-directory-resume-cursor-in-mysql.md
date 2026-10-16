# Store the per-directory resume cursor in MySQL

Status: accepted

vfm remembers where a user left off in a directory — the *resume cursor* — so that reopening the directory returns them there. A cursor is per user, per directory, and per sort order: reading a directory under name-ascending and under time-descending leaves two independent positions, each pointing at the last file read under that order. It lives today in a Redis hash `vfm:dir:lastpage:<userNo>`, fielded by dirKey alone, so one position exists per directory and a change of sort overwrites the other order's position; each field is a `{fileKey, time}` JSON blob with no TTL. We move it to a MySQL table, one row per (user, directory, sort), upserted as the user reads. The Redis hash keeps its existing entries and receives no further writes; it is read once, for a comic directory whose position exists only there, and the position found is carried into the table.

The same rows back the "directory browse history" list, which shows the directories a user read in most recently. That list becomes one indexed range scan over the table, replacing `HGETALL` plus a sort in the application plus a separate MySQL pass to enrich names and thumbnails. MySQL is also where `file_info`, the source of truth the cursor points into, already lives, so the cursor survives Redis eviction and restarts with it.

## Considered Options

- **Keep the Redis hash.** Cheapest to write and needs no migration, but offers no query for the recency list and silently loses cursors on eviction or restart.
- **Redis with a composite hash field (directory + sort).** Holds the per-sort cursors, but leaves the recency list as `HGETALL` + in-application sort + a separate enrichment pass, and the cursors still vanish on eviction.

## Consequences

- Cursor reads and writes now hit MySQL rather than Redis.
