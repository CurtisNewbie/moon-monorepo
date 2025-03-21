# mini-fstore

A mini file storage service. mini-fstore internally uses [miso](https://github.com/curtisnewbie/miso).

> **_This project is part of a monorepo ([https://github.com/CurtisNewbie/moon-monorepo](https://github.com/CurtisNewbie/moon-monorepo))._**

## Requirements

- MySQL
- Redis
- Consul
- RabbitMQ
- ffmpeg

## Configuration

For more configuration, see [miso](https://github.com/curtisnewbie/miso).

| Property                           | Description                                                                                                                                                                                                                               | Default Value |
| ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| fstore.storage.dir                 | Storage Directory                                                                                                                                                                                                                         | ./storage     |
| fstore.trash.dir                   | Trash Directory                                                                                                                                                                                                                           | ./trash       |
| fstore.tmp.dir                     | Temporary directory                                                                                                                                                                                                                       | /tmp          |
| fstore.pdelete.strategy            | Strategy used to 'physically' delete files, there are two types of strategies available: direct / trash. When using 'direct' strategy, files are deleted directly. When using 'trash' strategy, files are moved into the trash directory. | trash         |
| fstore.backup.enabled              | Enable endpoints for mini-fstore file backup, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup).                                                                                                                         | false         |
| fstore.backup.secret               | Secret for backup endpoints authorization, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup).                                                                                                                            |               |
| task.sanitize-storage-task.dry-run | Enable dry-run mode for StanitizeStorageTask                                                                                                                                                                                              | false         |

## Prometheus Metrics

- `mini_fstore_generate_file_key_duration`: histogram for monitoring time duration of random file key generation.
- `mini_fstore_generate_img_thumbnail_duration`: histogram for monitoring time duration of image thumbnail generation.
- `mini_fstore_generate_video_thumbnail_duration`: histogram for monitoring time duration of video thumbnail generation.

## Media Streming

The `/file/stream` endpoint can be used for media streaming.

```html
<body>
  <video controls>
    <source
      src="http://localhost:8084/file/stream?key=0fR1H1O0t8xQZjPzbGz4lRx%2FbPacIg"
      type="video/mp4"
    />
    Yo something is wrong
  </video>
</body>
```

## Limitation

Currently, mini-fstore nodes must all share the same database and the same storage devices. Some sort of distributed file system can be used and shared among all mini-fstore nodes if necessary.

## Docs

- [API Doc](./doc/api.md)
- [Workflows](./doc/workflow.md)

## Tools

- File Backup Tools: [fstore_backup](https://github.com/CurtisNewbie/fstore_backup).

## Maintenance

mini-fstore automatically detects duplicate uploads by comparing size and sha1 checksum. If duplicate file is detected, these files are _symbolically_ linked to the same file previously uploaded. This can massively reduce file storage, but multiple file records (multiple file_ids) can all point to a single file.

Whenever a file is marked logically deleted, the file is not truely deleted. In order to cleanup the storage for the deleted files including those that are possibly symbolically linked, you have to use the following endpoint to trigger the maintenance process. During the maintenance, uploading files is rejected.

```sh
curl -X POST http://localhost:8084/maintenance/remove-deleted
```

mini-fstore also provides maintenance endpoint that sanitize storage directory. Sometimes files are uploaded to storage directory, but are somehow not saved in database. These <i>dangling</i> files are handled by this endpoint.

```sh
curl -X POST http://localhost:8084/maintenance/sanitize-storage
```

To compute sha1 for previously uploaded files, use the following maintenance endpoint to trigger a compensation.

```sh
curl -X POST 'http://localhost:8084/maintenance/compute-checksum'
```

## Update
