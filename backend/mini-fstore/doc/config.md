# Configurations

For more configuration, see [github.com/curtisnewbie/miso](https://github.com/CurtisNewbie/miso/blob/main/doc/config.md).

## mini-fstore configuration

| property                           | description                                                                                                                                                                                                                               | default value |
| ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| fstore.storage.dir                 | Storage Directory                                                                                                                                                                                                                         | ./storage     |
| fstore.trash.dir                   | Trash Directory                                                                                                                                                                                                                           | ./trash       |
| fstore.tmp.dir                     | Temporary directory                                                                                                                                                                                                                       | /tmp          |
| fstore.pdelete.strategy            | Strategy used to 'physically' delete files, there are two types of strategies available: direct / trash. When using 'direct' strategy, files are deleted directly. When using 'trash' strategy, files are moved into the trash directory. | `"trash"`     |
| fstore.backup.enabled              | Enable endpoints for mini-fstore file backup, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup).                                                                                                                         | false         |
| fstore.backup.secret               | Secret for backup endpoints authorization, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup).                                                                                                                            |               |
| task.sanitize-storage-task.dry-run | Enable dry-run mode for StanitizeStorageTask                                                                                                                                                                                              | false         |