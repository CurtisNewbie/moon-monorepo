# Configurations

For more configuration, see [github.com/curtisnewbie/miso](https://github.com/CurtisNewbie/miso/blob/main/doc/config.md).

## General Configuration

| property      | description                | default value |
| ------------- | -------------------------- | ------------- |
| vfm.site.host | Externally accessible host |               |

## VFM Configuration

| property                 | description                                                                                                  | default value |
| ------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------- |
| vfm.temp-path            | Temporary file path for bootmarks files                                                                      | "/tmp/vfm"    |
| vfm.unpack.max-depth     | Maximum depth of an archive's folder tree when unpacking, counted in path segments from the target directory | 32            |
| vfm.zip-tree.max-depth   | Maximum path segment depth included in a zip tree preview                                                    | 32            |
| vfm.zip-tree.max-entries | Maximum tree nodes represented in a zip tree preview                                                         | 10000         |