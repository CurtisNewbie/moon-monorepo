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

-- for migration
insert into user_vault.path select * from goauth.path;
insert into user_vault.path_resource select * from goauth.path_resource;
insert into user_vault.resource select * from goauth.resource;
insert into user_vault.role select * from goauth.role;