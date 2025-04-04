# mini-fstore

A mini file storage service. mini-fstore internally uses [miso](https://github.com/curtisnewbie/miso).

> **_This project is part of a monorepo ([https://github.com/CurtisNewbie/moon-monorepo](https://github.com/CurtisNewbie/moon-monorepo))._**

## Requirements

- ffmpeg

## Limitation

Currently, mini-fstore nodes must all share the same database and the same storage devices. Some sort of distributed file system can be used and shared among all mini-fstore nodes if necessary.

## Docs

- [Configuration](./doc/config.md)
- [Metrics](./doc/metrics.md)
- [API Doc](./doc/api.md)
- [Maintenance](./doc/maintenance.md)
- [Media Streaming](./doc/streaming.md)
- [Workflows](./doc/workflow.md)

## Tools

- File Backup Tools: [fstore_backup](https://github.com/CurtisNewbie/fstore_backup).
