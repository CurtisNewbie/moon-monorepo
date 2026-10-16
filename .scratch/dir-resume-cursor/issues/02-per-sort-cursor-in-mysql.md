# 02: Per-sort resume cursor in MySQL, end to end

**What to build:** Every directory, comic or not, remembers the file a user last read under each sort order, and reopening a directory offers to jump to the page that file sits on under the sort in force. The memory moves out of Redis into a MySQL table.

A cursor row, per user, directory and sort order:

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

Cursors are soft-deleted and read only where `deleted` is zero. Because the lookup columns a cursor is keyed by cannot carry a unique constraint under soft delete, writes are serialized by a Redis lock taken on (user, directory, sort), with a read-then-insert-or-update under the lock.

The existing per-directory last-page endpoints carry the requested sort on both the record and the lookup. The effective sort is resolved on the server and is the same on both paths, so a comic directory — whose client sends its dropdown default — records and reads under one key. Recording is limited to the directory's owner and to directories that exist; the comic-only gate is removed so every directory records. The page is never stored: the client locates the recorded file through the file-position endpoint under the sort in force and the current page size, then shows the existing Continue Browsing dialog.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] A non-comic directory records a cursor as its files are read, and reopening it offers Continue Browsing and jumps to the recorded file's page.
- [ ] A comic directory resumes exactly as it does today.
- [ ] Reading a directory under a second sort leaves the first sort's cursor intact; each sort resumes independently.
- [ ] The offered page matches where the recorded file sits under the sort in force and the current page size; neither page nor page size is stored.
- [ ] Concurrent recordings for the same (user, directory, sort) leave exactly one cursor row.
- [ ] A recorded file that no longer exists opens at page 1 without error.
- [ ] No cursor is recorded for a directory the user does not own, or for a directory the user has not read three distinct consecutive files in.
- [ ] The table appears in the vfm schema snapshot and in the in-flight release increment.
