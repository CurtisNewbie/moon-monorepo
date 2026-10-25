-- Drop the fulltext index on file name
-- The file listing now matches name substrings per term, so this index serves no query
ALTER TABLE vfm.file_info DROP INDEX `name_idx`;
