alter table file_info
    modify column user_group int NOT NULL default 1 COMMENT 'the group that the file belongs to, 0-public, 1-private',
    add column uploader_no varchar(32) NOT NULL DEFAULT '' COMMENT 'user no of uploader',
    add key uploader_no_idx (uploader_no);
