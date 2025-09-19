CREATE TABLE user_vault.note (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
    `record_id` varchar(40) NOT NULL DEFAULT '' COMMENT 'record_id',
    `title` text NULL DEFAULT NULL COMMENT 'title',
    `content` mediumtext NULL DEFAULT NULL COMMENT 'content' ,
    `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
    `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
    `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
    `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
    `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted',
    PRIMARY KEY (`id`),
    KEY user_no_idx (user_no, deleted),
    KEY record_id_idx (record_id, deleted),
    FULLTEXT `title_content_idx` (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='notes';

