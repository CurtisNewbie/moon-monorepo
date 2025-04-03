# Updates

- Since v0.0.4, `vfm` relies on `evnet-pump` to listen to the binlog events. Whenever a new `file_info` record is inserted, the `event-pump` sends MQ to `vfm`, which triggers the image compression workflow if the file is an image.
- Since v0.0.8, Users can only share files using `vfolder`, field `file_info.user_group` and table `file_sharing` are deprecated.
- Since v0.1.3, [fantahsea](https://github.com/curtisnewbie/fantahsea) has been merged into vfm codebase, see [Fantahsea Migration](./doc/fantahsea-migration.md).
- Since v0.1.17, file tag functionality is removed.