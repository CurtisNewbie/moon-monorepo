# Configurations

For more configuration, see [github.com/curtisnewbie/miso](https://github.com/CurtisNewbie/miso/blob/main/doc/config.md).

## Logbot Configuration

| property                        | description                                             | default value |
| ------------------------------- | ------------------------------------------------------- | ------------- |
| logbot.node                     | logbot node name                                        | default       |
| logbot.watch                    | (`slice of watch object`) logbot watch configuration    |               |
| logbot.watch.app                | (`watch object`) app name                               |               |
| logbot.watch.file               | (`watch object`) path of the log file                   |               |
| logbot.watch.type               | (`watch object`) type of log pattern `[ 'go', 'java' ]` |               |
| logbot.remove-history-error-log | enable task to remove error logs reported 7 days ago    | false         |
| log.pattern                     | (`slice of string`) log pattern supported (regexp)      |               |
| log.merged-file-name            | merged log filename                                     |               |