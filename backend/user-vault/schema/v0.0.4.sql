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

alter table user_vault.user add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  change column `create_by` `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  change column `update_by` `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  change column `is_del` `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted';


alter table user_vault.user_key add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  change column `create_by` `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  change column `update_by` `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  change column `is_del` `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted';

alter table user_vault.user_key comment 'user key';

alter table user_vault.access_log add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  add column `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  add column `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  add column `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  add column `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  add column `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted',
  add key username_idx (username, deleted);

alter table user_vault.path add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  change column `create_time` `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  change column `update_time` `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  change column `create_by` `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  change column `update_by` `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  change column `is_del` `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted';


alter table user_vault.path_resource add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  change column `create_time` `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  change column `update_time` `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  change column `create_by` `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  change column `update_by` `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  change column `is_del` `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted';

alter table user_vault.resource add column `trace_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'trace_id',
  change column `create_time` `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  change column `update_time` `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  change column `create_by` `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'created by',
  change column `update_by` `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'updated by',
  change column `is_del` `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'record deleted';