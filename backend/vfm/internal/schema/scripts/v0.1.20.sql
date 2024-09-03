alter table file_info drop column uploader_id,
    modify column is_logic_deleted tinyint NOT NULL DEFAULT '0' COMMENT 'whether the file is logically deleted, 0-normal, 1-deleted',
    modify column is_physic_deleted tinyint NOT NULL DEFAULT '0' COMMENT 'whether the file is physically deleted, 0-normal, 1-deleted',
    drop column fs_group_id,
    add column hidden tinyint not null default '0' comment 'whether the file is hidden';

CREATE TABLE IF NOT EXISTS versioned_file (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'primary key',
    name varchar(255) NOT NULL COMMENT 'name of the file',
    ver_file_id VARCHAR(32) NOT NULL COMMENT 'versioned file id',
    file_key VARCHAR(64) NOT NULL COMMENT 'file_info key',
    size_in_bytes BIGINT NOT NULL COMMENT 'size in bytes',
    uploader_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'uploader user_no',
    uploader_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'uploader name',
    upload_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'upload time',
    delete_time DATETIME DEFAULT NULL COMMENT 'when the file is logically deleted',
    deleted TINYINT(4) NOT NULL DEFAULT '0' COMMENT 'whether the file is deleted, 0-false, 1-true',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
    created_by VARCHAR(255) NOT NULL DEFAULT '' comment 'created by',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
    updated_by VARCHAR(255) NOT NULL DEFAULT '' comment 'updated by',
    UNIQUE KEY ver_file_id_uk (ver_file_id),
    KEY file_key_idx (file_key),
    KEY uploader_no_idx (uploader_no),
    FULLTEXT KEY name_idx (name)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='Versioned File';

CREATE TABLE IF NOT EXISTS versioned_file_log (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'primary key',
    ver_file_id VARCHAR(32) NOT NULL COMMENT 'versioned file id',
    file_key VARCHAR(64) NOT NULL COMMENT 'file_info key',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
    created_by VARCHAR(255) NOT NULL DEFAULT '' comment 'created by',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
    updated_by VARCHAR(255) NOT NULL DEFAULT '' comment 'updated by',
    KEY ver_file_id_idx (ver_file_id)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='Versioned File Log';

CREATE TABLE IF NOT EXISTS bookmark (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `icon` text COMMENT 'icon encoded blob',
  `name` varchar(512) NOT NULL DEFAULT '' COMMENT 'bookmark name',
  `href` varchar(1024) NOT NULL DEFAULT '' COMMENT 'bookmark href',
  `md5` varchar(32) not null default '' comment 'md5',
  PRIMARY KEY (`id`),
  KEY `idx_name` (name),
  UNIQUE `uk_user_no_md5` (user_no, md5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Bookmark';

CREATE TABLE IF NOT EXISTS bookmark_blacklist (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `icon` text COMMENT 'icon encoded blob',
  `name` varchar(512) NOT NULL DEFAULT '' COMMENT 'bookmark name',
  `href` varchar(1024) NOT NULL DEFAULT '' COMMENT 'bookmark href',
  `md5` varchar(32) not null default '' comment 'md5',
  PRIMARY KEY (`id`),
  UNIQUE `uk_user_no_md5` (user_no, md5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Blacklisted Bookmark';