CREATE DATABASE IF NOT EXISTS logbot;

CREATE TABLE `error_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `node` varchar(25) NOT NULL COMMENT 'node name',
  `app` varchar(25) NOT NULL COMMENT 'app name',
  `caller` varchar(255) NOT NULL COMMENT 'caller name',
  `trace_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `span_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `err_msg` text COMMENT 'error msg',
  `rtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'report time',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`),
  KEY `idx_rtime` (`rtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Application Error Log';

