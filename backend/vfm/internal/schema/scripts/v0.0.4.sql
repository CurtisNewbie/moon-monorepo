alter table user_vfolder add column username varchar(50) default '' comment 'username';

-- For the new `username` column
-- update fileserver.user_vfolder uv left join authserver.user u on u.user_no = uv.user_no set uv.username = u.username where uv.username = '';