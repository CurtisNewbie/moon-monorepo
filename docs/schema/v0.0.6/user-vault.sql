-- API path descriptions can outgrow 255 chars (e.g. fstore's /fstore/file/bundle)
ALTER TABLE user_vault.path MODIFY COLUMN `desc` varchar(1024) NOT NULL DEFAULT '' COMMENT 'description';
