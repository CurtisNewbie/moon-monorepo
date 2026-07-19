-- Add seq_key column for custom drag-and-drop file ordering
-- seq_key stores fractional-indexing key as a base-95 encoded string
-- Empty string (default) means file is not in custom order
-- Use binary collation for case-sensitive base-62 key comparison
ALTER TABLE vfm.file_info ADD COLUMN seq_key VARCHAR(50) NOT NULL DEFAULT '' COLLATE utf8mb4_bin COMMENT 'fractional-indexing key for custom drag-and-drop ordering';
ALTER TABLE vfm.file_info ADD KEY `idx_parent_seq_key` (`parent_file`, `seq_key`);

