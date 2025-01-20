# gatekeeper

Simple custom-made gateway written in Go. This project is internally backed by [curtisnewbie/miso](https://github.com/curtisnewbie/miso).

> **_This project is part of a monorepo ([https://github.com/CurtisNewbie/moon-monorepo](https://github.com/CurtisNewbie/moon-monorepo))._**

## Dependencies

- Consul
- [github.com/curtisnewbie/user-vault](https://github.com/curtisnewbie/user-vault)

## Configuration

See [miso](https://github.com/curtisnewbie/miso) for more about configuration.

| Property                           | Description                                                                                                                                                                                                   | Default Value |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| gatekeeper.timer.path.excl         | slice of paths that are not measured by prometheus timer                                                                                                                                                      |               |
| gatekeeper.whitelist.path.patterns | slice of path patterns that do not require authorization and authentication                                                                                                                                   |               |
| gatekeeper.overwrite-remote-ip     | always overwrite remote ip address in `x-forwarded-for` header (by default, there should be a nginx sitting right before the gatekeeper as a reverse proxy, this the default value for this setting is false) | false         |
