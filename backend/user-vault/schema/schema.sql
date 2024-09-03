create database if not exists user_vault;

CREATE TABLE IF NOT EXISTS user_vault.user (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `username` varchar(50) NOT NULL COMMENT 'username',
  `password` varchar(255) NOT NULL COMMENT 'password in hash',
  `salt` varchar(10) NOT NULL COMMENT 'salt',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `is_disabled` int NOT NULL DEFAULT '0' COMMENT 'whether the user is disabled, 0-normal, 1-disabled',
  `review_status` varchar(25) NOT NULL COMMENT 'Review Status',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `user_no` (`user_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User';

CREATE TABLE IF NOT EXISTS user_vault.user_key (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_id` int unsigned NOT NULL COMMENT 'user.id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'name of the key',
  `secret_key` varchar(255) NOT NULL COMMENT 'secret key',
  `expiration_time` datetime NOT NULL COMMENT 'when the key is expired',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no',
  PRIMARY KEY (`id`),
  UNIQUE KEY `secret_key` (`secret_key`),
  KEY `user_no_idx` (`user_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='user''s key'


CREATE TABLE user_vault.access_log (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `access_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the user signed in',
  `ip_address` varchar(255) NOT NULL COMMENT 'ip address',
  `username` varchar(255) NOT NULL COMMENT 'username',
  `user_id` int unsigned NOT NULL COMMENT 'primary key of user',
  `url` varchar(255) DEFAULT '' COMMENT 'request url',
  `user_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'User Agent',
  `success` tinyint(1) DEFAULT '1' COMMENT 'login was successful',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='access log';

CREATE TABLE IF NOT EXISTS user_vault.path (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `pgroup` varchar(20) NOT NULL DEFAULT '' COMMENT 'path group',
  `path_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'path no',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT 'description',
  `method` varchar(10) NOT NULL DEFAULT ''  COMMENT 'http method',
  `url` varchar(128) NOT NULL DEFAULT '' COMMENT 'path url',
  `ptype` varchar(10) NOT NULL DEFAULT '' COMMENT 'path type: PROTECTED, PUBLIC',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `path_no` (`path_no`)
) ENGINE=InnoDB COMMENT='Paths';

CREATE TABLE IF NOT EXISTS user_vault.path_resource (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `path_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'path no',
  `res_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource code',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY (`path_no`, `res_code`)
) ENGINE=InnoDB COMMENT='Path Resource';

CREATE TABLE IF NOT EXISTS user_vault.resource (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `code` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource code',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource name',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `code` (`code`)
) ENGINE=InnoDB COMMENT='Resources';

CREATE TABLE IF NOT EXISTS user_vault.role_resource (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `res_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource code',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB COMMENT='Role resources';

CREATE TABLE IF NOT EXISTS user_vault.role (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'name of role',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB COMMENT='Roles';

CREATE TABLE IF NOT EXISTS notification (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
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
  KEY `user_no_status_idx` (`user_no`, `status`),
  UNIQUE KEY `notifi_no_uk` (`notifi_no`)
) ENGINE=InnoDB COMMENT='Platform Notification';

-- default one for administrator, with this role, all paths can be accessed
INSERT INTO user_vault.role(role_no, name) VALUES ('role_554107924873216177918', 'Super Administrator');