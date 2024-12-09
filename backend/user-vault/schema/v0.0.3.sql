alter table user_vault.path_resource modify column res_code varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code';

alter table user_vault.resource modify column code varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code';

alter table user_vault.role_resource modify column res_code varchar(128) NOT NULL DEFAULT '' COMMENT 'resource code';