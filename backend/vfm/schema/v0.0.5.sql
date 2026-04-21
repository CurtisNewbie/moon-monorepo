-- Add index for optimizing batch directory thumbnail queries
-- This index improves performance of finding the newest child file (by id) for each directory
ALTER TABLE vfm.file_info ADD KEY `idx_parent_thumbnail` (`parent_file`,`id` DESC);