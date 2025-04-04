# gatekeeper

Simple custom-made gateway written in Go. This project is internally backed by [curtisnewbie/miso](https://github.com/curtisnewbie/miso).

> **_This project is part of a monorepo ([https://github.com/CurtisNewbie/moon-monorepo](https://github.com/CurtisNewbie/moon-monorepo))._**

## Configuration

See [miso](https://github.com/curtisnewbie/miso) for more about configuration.

| Property                           | Description                                                                                                                                                                                                   | Default Value |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| gatekeeper.timer.path.excl         | slice of paths that are not measured by prometheus timer                                                                                                                                                      |               |
| gatekeeper.whitelist.path.patterns | slice of path patterns that do not require authorization and authentication                                                                                                                                   |               |
| gatekeeper.overwrite-remote-ip     | always overwrite remote ip address in `x-forwarded-for` header (by default, there should be a nginx sitting right before the gatekeeper as a reverse proxy, this the default value for this setting is false) | false         |
| gatekeeper.proxy.pprof.bearer      | bearer authentication token for pprof endpoints (not just gatekeeper, but also servers behind it), mandatory for production mode; if bearer token is missing, all requests to `*/debug/pprof/*` are rejected  |

## Authentication and Authorization

gatekeeper relies on user-vault to do user authentication and authorization. user-vault manages all the user information, while gatekeepr acts as a gateway just passing the bearer tokens to user-vault (see user-vault's documentation for more details).

There are two ways to pass authentication tokens to gatekeeper:

1. Set the Bearer token to the `Authorization` request header.
2. Set the Bearer token to the `Gatekeeper_Authorization` cookie.

gatekeeper extracts the token from the request and passes the token to user-vault, if user-vault recognizes the token and authorizes the request, the request is proxied to backend servers behind gatekeeper.

## Proxy pprof requests with bearer authentication

1. Set value to property `gatekeeper.proxy.pprof.bearer`, the expected bearer token.
2. Enable proxied apps pprof endpoints, e.g.,

   ```yaml
   server:
     pprof:
       enabled: true
   ```

3. Use curl to retrieve pprof file, e.g.,

   ```sh
   tok="" # your bearer token
   sec="30"
   server="$1"
   out="/tmp/${server}_heap.pprof"

   curl https://$server/debug/pprof/heap?seconds=$sec -H "Authorization: Bearer $tok" -v -o $out \
   ```

4. Use go tool pprof to open the downloaded file:

   ```sh
   go tool pprof -http=: "$out"
   ```
