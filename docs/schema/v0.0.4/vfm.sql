ALTER TABLE vfm.gallery_image ADD KEY gallery_name_idx (gallery_no, name), DROP KEY gallery_no_idx;

ALTER TABLE vfm.file_info ADD KEY `parent_file_create_time_idx` (`parent_file`,`create_time`);