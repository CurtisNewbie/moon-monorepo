create database if not exists user_vault;

-- TODO update ddl
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `username` varchar(50) NOT NULL COMMENT 'username',
  `password` varchar(255) NOT NULL COMMENT 'password in hash',
  `salt` varchar(10) NOT NULL COMMENT 'salt',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `is_disabled` int(11) NOT NULL DEFAULT '0' COMMENT 'whether the user is disabled, 0-normal, 1-disabled',
  `review_status` varchar(25) NOT NULL COMMENT 'Review Status',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `user_no` (`user_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User';

CREATE TABLE `user_key` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_id` int(10) unsigned NOT NULL COMMENT 'user.id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'name of the key',
  `secret_key` varchar(255) NOT NULL COMMENT 'secret key',
  `expiration_time` datetime NOT NULL COMMENT 'when the key is expired',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
  PRIMARY KEY (`id`),
  UNIQUE KEY `secret_key` (`secret_key`),
  KEY `user_no_idx` (`user_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='user key';

CREATE TABLE `access_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `access_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the user signed in',
  `ip_address` varchar(255) NOT NULL COMMENT 'ip address',
  `username` varchar(255) NOT NULL COMMENT 'username',
  `user_id` int(10) unsigned NOT NULL COMMENT 'primary key of user',
  `url` varchar(255) DEFAULT '' COMMENT 'request url',
  `user_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'User Agent',
  `success` tinyint(1) DEFAULT '1' COMMENT 'login was successful',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='access log';

CREATE TABLE `path` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `pgroup` varchar(20) NOT NULL DEFAULT '' COMMENT 'path group',
  `path_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'path no',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT 'description',
  `method` varchar(10) NOT NULL DEFAULT '' COMMENT 'http method',
  `url` varchar(128) NOT NULL DEFAULT '' COMMENT 'path url',
  `ptype` varchar(10) NOT NULL DEFAULT '' COMMENT 'path type: PROTECTED, PUBLIC',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `path_no` (`path_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Paths';

CREATE TABLE `path_resource` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `path_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'path no',
  `res_code` varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `path_no` (`path_no`,`res_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Path Resource';

CREATE TABLE `resource` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `code` varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'resource name',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Resources';

CREATE TABLE `role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'name of role',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Roles';

CREATE TABLE `role_resource` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `res_code` varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Role resources';

CREATE TABLE `notification` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `notifi_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'notification no',
  `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT 'title',
  `message` varchar(1000) NOT NULL DEFAULT '' COMMENT 'message',
  `status` varchar(10) NOT NULL DEFAULT 'INIT' COMMENT 'Status: INIT, OPENED',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  PRIMARY KEY (`id`),
  UNIQUE KEY `notifi_no_uk` (`notifi_no`),
  KEY `user_no_status_idx` (`user_no`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Platform Notification';

CREATE TABLE `site_password` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `record_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'record unique id',
  `site` varchar(64) NOT NULL DEFAULT '' COMMENT 'site',
  `alias` varchar(64) NOT NULL DEFAULT '' COMMENT 'alias',
  `username` varchar(50) NOT NULL COMMENT 'username',
  `password` varchar(255) NOT NULL COMMENT 'site password encrypted using user login password',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `record_id_uk` (`record_id`),
  KEY `user_alias_idx` (`user_no`,`alias`),
  KEY `user_site_idx` (`user_no`,`site`),
  KEY `user_username_idx` (`user_no`,`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Personal passwords for different sites';

CREATE TABLE note (
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

-- default one for administrator, with this role, all paths can be accessed
INSERT INTO user_vault.role(role_no, name) VALUES ('role_super_admin', 'Super Administrator');