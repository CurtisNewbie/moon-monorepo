# Configurations

For more configuration, see [github.com/curtisnewbie/miso](https://github.com/CurtisNewbie/miso/blob/main/doc/config.md).

## Gatekeeper Configuration

| property                           | description                                                                                                                                                                                                   | default value |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| gatekeeper.timer.path.excl         | slice of paths that are not measured by prometheus timer                                                                                                                                                      |               |
| gatekeeper.whitelist.path.patterns | slice of path patterns that do not require authorization and authentication                                                                                                                                   |               |
| gatekeeper.overwrite-remote-ip     | always overwrite remote ip address in `x-forwarded-for` header (by default, there should be a nginx sitting right before the gatekeeper as a reverse proxy, this the default value for this setting is false) | false         |
| gatekeeper.proxy.pprof.bearer      | bearer authentication token for pprof endpoints (not just gatekeeper, but also servers behind it), mandatory for production mode; if bearer token is missing, all requests to `*/debug/pprof/*` are rejected  |               |