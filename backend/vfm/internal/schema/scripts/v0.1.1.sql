alter table file_info add column `sensitive_mode` varchar(1) NOT NULL DEFAULT 'N' COMMENT 'sensitive file, Y/N';