CREATE TABLE IF NOT EXISTS user_vault.site_password (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
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
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `record_id_uk` (`record_id`),
  KEY `user_alias_idx` (`user_no`, `alias`),
  KEY `user_site_idx` (`user_no`, `site`),
  KEY `user_username_idx` (`user_no`, `username`)
) ENGINE=InnoDB COMMENT='Personal passwords for different sites';