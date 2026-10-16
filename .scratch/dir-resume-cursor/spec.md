# Directory Resume Cursor

Status: ready-for-agent

## Problem Statement

A user reading files in a directory has no reliable way back to where they stopped.

Today vfm records a "last page" only for directories marked as comic; any other directory forgets the user's place entirely. Even where it records, there is exactly one position per directory, so reading the same directory under a different sort order overwrites the other order's position — and a page number is meaningless under a sort it was not recorded under. Reopening a directory always starts at page 1 unless the user goes through the history list.

The recorded position lives only in Redis, with no TTL, no query, and no durability: it disappears on eviction or restart, and there is nothing for an operator to inspect or clean.

## Solution

Every directory remembers a Resume Cursor: the file the user last read, kept per user, per directory, and per sort order.

When a user opens a directory, vfm resolves the sort order in force for that directory, looks up the cursor for that sort, and — if one exists — offers to jump to the page that file sits on under that sort, using the existing "Continue Browsing" dialog. The user can accept or start from the top.

The cursor moves from Redis to a MySQL table so it is durable, queryable, and soft-deletable. Cursors already in Redis keep working: they are read once as a fallback and written back into the table.

## User Stories

1. As a moon user, I want every directory I read to remember the last file I read, so that reopening it returns me to where I stopped.
2. As a moon user, I want that memory kept separately per sort order, so that reading a directory under a second sort does not destroy the position I had under the first.
3. As a moon user, I want the position offered to me only under the sort order it was recorded under, so the page number I am offered is meaningful.
4. As a moon user, I want to be asked before jumping, so that I can choose to start from the top instead.
5. As a moon user, I want the dialog to state the page number, so that I know where I would land.
6. As a moon user, I want to jump to a page whose file list is already ordered the way it was when I left, so that the page looks the way I expect.
7. As a moon user reading a comic directory, I want the same resume behavior I have today, so that nothing changes for me.
8. As a moon user who only glanced at a directory (fewer than three distinct files in a row), I want no cursor recorded, so that an accidental click does not move my place.
9. As a moon user, I want the resume position I already had before this change to keep working, so that I do not lose where I was.
10. As a moon user, I want the position I already had to be carried into the new store automatically, so that it survives the change even if I never revisit the directory.
11. As a moon user, I want the Directory Browse History list to show each directory once, from its most recent visit, so that it is not cluttered by one entry per sort.
12. As a moon user, I want resuming from the Directory Browse History list to land on the file I last read there under the sort I last used, so that the list's Resume action is faithful.
13. As a moon user, I want a directory I deleted to disappear from the Directory Browse History list.
14. As a moon user, I want my cursor updated as I read, so that the position is current when I return.
15. As a moon user browsing quickly, I do not want my reading to flood the database with writes.
16. As a moon user, I want the resume to locate the correct file regardless of the page size in effect, so that a layout or device change does not send me to the wrong place.
17. As a moon user, I want the sort order I chose for a directory to keep opening that directory, unaffected by this feature.
18. As a moon user, I want a directory whose recorded file no longer exists to open at page 1 without an error.
19. As a moon user, I want my cursor visible to me alone, so that no other user's reading position is ever exposed to me.
20. As an operator, I want cursor rows to live in a table I can query, so that I can inspect and clean them.
21. As an operator, I want to discard a directory's cursor by marking it deleted, so that the removal is auditable rather than destructive.
22. As an operator releasing a version, I want the new table recorded in the schema snapshot, the release increment, and the change log, so that the schema stays reproducible.

## Implementation Decisions

**Storage.** A new MySQL table, `vfm.dir_resume_cursor`, one row per user, directory, and sort order. A cursor is soft-deletable, so the lookup columns carry a plain index, not a unique key; concurrent writes are serialized by a Redis lock keyed on the same tuple, and the writer reads-then-inserts-or-updates under that lock.

```sql
CREATE TABLE vfm.dir_resume_cursor (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `dir_key` varchar(64) NOT NULL COMMENT 'directory file_key (file_info.uuid)',
  `order_by` varchar(16) NOT NULL COMMENT 'sort field: name, time, custom',
  `order_dir` varchar(4) NOT NULL COMMENT 'sort direction: asc, desc',
  `file_key` varchar(64) NOT NULL COMMENT 'last-read file uuid under this sort',
  PRIMARY KEY (`id`),
  KEY `user_dir_sort_idx` (`user_no`,`dir_key`,`order_by`,`order_dir`),
  KEY `user_updated_idx` (`user_no`,`updated_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Per-directory resume cursor, one row per user, directory and sort order';
```

A row's columns and their roles:

| Column | Role |
| --- | --- |
| `user_no` | owner of the cursor; the same user is the only one who reads it |
| `dir_key` | the directory, as the `file_info.uuid` of a directory row |
| `order_by` | the sort field the cursor was recorded under |
| `order_dir` | the sort direction, always resolved to `asc` or `desc` |
| `file_key` | the last file read under that sort |
| `deleted` | soft delete; the cursor is only ever read where it is zero |

**Effective sort.** A single server-side resolution decides the sort order from the requested values: a comic directory is forced to `name` ascending, and an absent direction is resolved to `asc` for `name`/`custom` and `desc` for time. The same resolved tuple is used when recording a cursor and when looking one up, so a comic directory — whose client sends its dropdown default — records and reads under the same key.

**Cursor contract.** The existing per-directory last-page endpoints are extended to carry the sort: recording takes the directory, the file, and the requested sort; the lookup takes the directory and the requested sort and returns the recorded file (empty when none). The lookup returns the raw file, not a page.

**Page derivation.** The page is never stored. On resume the client locates the recorded file through the existing file-position endpoint, passing the resolved sort and the current page size, and obtains the page under that sort. Page size is never stored, so the offered page reflects the current page size while the located file is unaffected by it.

**Recording.** Every directory records a cursor — the comic-only gate is removed. Recording happens from the client's existing debounced pipeline; no coalescing is added on the server. Recording remains limited to the directory's owner and to directories that exist.

**Legacy compatibility.** On a lookup miss, and only for a comic directory (the only kind with entries under the old key), the old Redis last-page hash is consulted; a hit is written into the table and then returned. Existing Redis entries are otherwise left in place and never written to; the Redis hash is no longer a write target.

**Directory Browse History.** The list returns one entry per directory, taken from its most recently updated cursor row, and excludes directories marked deleted. Each entry carries the sort order and file of that row, and resuming from it restores that sort and file rather than the fixed `name` sort it uses today.

**Schema records.** The table is added to the vfm schema snapshot and appended to the in-flight release increment, both for the same release; the release's change-log entry already lists the vfm increment.

## Testing Decisions

Tests assert external behavior only: given a directory and a sort, what a cursor lookup returns and how it changes as files are read and directories are opened. No test reaches into the table's columns or asserts on the locking mechanism.

- Backend seam: the vfm service layer, exercised in-process against a real MySQL and Redis, as the existing service tests in that package already do. It covers effective-sort resolution (comic forced, direction defaulted), one row per sort, that recording under a second sort leaves the first intact, the locked read-then-write under concurrent recording, the legacy fallback with write-back, the list returning one entry per directory excluding deleted directories, and that a missing recorded file yields no resume rather than an error.
- Frontend seam: the browse-history recorder service, unit-tested with an HTTP controller mock as its existing spec does. It covers that recording and lookup now carry the directory and the sort, and that the recorder exposes a resolved page rather than a stored one.

Prior art: the package's existing service tests for the backend seam; the recorder's existing spec for the frontend seam.

## Out of Scope

- The comic directory's forced name-ascending ordering, in both listing and position calculation. It stays as it is; comic directories simply keep resuming under that order.
- Moving the per-directory Order-By Preference out of Redis. This feature reads it, it does not relocate it.
- Persisting page size, or making it configurable per user or directory.
- Server-side coalescing of cursor writes.
- Any change to what marks a directory as comic, or to the recording of the global per-file browse history.

## Further Notes

- The cursor's shape and its move to MySQL are recorded in ADR 0002; `Resume Cursor` and `Order-By Preference` are defined in `GLOSSARY.md`.
- There is no gate on either side of the cursor: every directory records and every directory reads. A cursor therefore cannot be orphaned by a flag change on its directory.
