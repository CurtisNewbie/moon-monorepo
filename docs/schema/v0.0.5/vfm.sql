-- Add index for optimizing batch directory thumbnail queries
-- This index improves performance of finding the newest child file (by id) for each directory
ALTER TABLE vfm.file_info ADD KEY `idx_parent_thumbnail` (`parent_file`,`id` DESC);
ALTER TABLE vfm.file_info ADD COLUMN is_comic TINYINT NOT NULL DEFAULT 0 COMMENT 'marks comic directory';
ALTER TABLE vfm.file_info ADD KEY `idx_parent_name` (`parent_file`, `name`);

-- Add seq_key column for custom drag-and-drop file ordering
-- seq_key stores fractional-indexing key as a base-95 encoded string
-- Empty string (default) means file is not in custom order
-- Use binary collation for case-sensitive base-62 key comparison
ALTER TABLE vfm.file_info ADD COLUMN seq_key VARCHAR(50) NOT NULL DEFAULT '' COLLATE utf8mb4_bin COMMENT 'fractional-indexing key for custom drag-and-drop ordering';
ALTER TABLE vfm.file_info ADD KEY `idx_parent_seq_key` (`parent_file`, `seq_key`);

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
