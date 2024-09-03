use user_vault;

alter table user_key add column `user_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'user no', add index user_no_idx (user_no);
update user_key uk left join user u on uk.user_id = u.id set uk.user_no = u.user_no;