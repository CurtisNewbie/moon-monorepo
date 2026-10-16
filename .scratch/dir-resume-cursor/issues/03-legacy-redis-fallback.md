# 03: Legacy Redis cursors carried forward

**What to build:** A user keeps the reading position they had before this change. When a directory's cursor is looked up and the table holds no row for it, and the directory is comic (the only kind the old store ever held), the old Redis last-page entry is consulted; if one is found, it is used as the cursor and written into the table, so later reads are served from the table.

**Blocked by:** 02

**Status:** ready-for-agent

- [ ] A comic directory whose position exists only in the old Redis store still resumes, using that entry's file.
- [ ] On that first read the position is written into the table, and subsequent reads are served from the table without touching Redis.
- [ ] A directory with neither a table row nor a Redis entry opens at page 1 without error.
- [ ] The Redis store receives no new writes.
