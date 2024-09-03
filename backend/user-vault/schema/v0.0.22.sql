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