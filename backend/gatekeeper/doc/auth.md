# Authentication and Authorization

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