USE vfm;

-- warn: this will rebuilt the table, may take a while
ALTER TABLE file_info
    DROP INDEX name,
    ADD FULLTEXT name_idx (name);

ALTER TABLE user_vfolder
    DROP INDEX user_folder_uk,
    ADD INDEX user_folder_idx (user_no, folder_no);