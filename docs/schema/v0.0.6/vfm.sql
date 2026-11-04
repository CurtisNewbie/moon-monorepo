-- Drop the fulltext index on file name
-- The file listing now matches name substrings per term, so this index serves no query
ALTER TABLE vfm.file_info DROP INDEX `name_idx`;

-- Index gallery_image by file_key alone
-- Lookups by file_key (e.g. removing an image by its file key) don't use the
-- existing (gallery_no, file_key) unique key
ALTER TABLE vfm.gallery_image ADD INDEX `file_key_idx` (`file_key`);
