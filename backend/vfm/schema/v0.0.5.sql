-- Add index for optimizing batch directory thumbnail queries
-- This index improves performance of finding the newest child file (by id) for each directory
ALTER TABLE vfm.file_info ADD KEY `idx_parent_thumbnail` (`parent_file`,`id` DESC);
ALTER TABLE vfm.file_info ADD COLUMN is_comic TINYINT NOT NULL DEFAULT 0 COMMENT 'marks comic directory';
ALTER TABLE vfm.file_info ADD KEY `idx_parent_name` (`parent_file`, `name`);