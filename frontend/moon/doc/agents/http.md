# HTTP interceptors

Two interceptors are registered in `src/app/app.module.ts`, in this order:

| Order | Interceptor       | Role                                                        |
| ----- | ----------------- | ----------------------------------------------------------- |
| 1     | `ErrorInterceptor` | Global error response handling (401, 403, 504, unknown)     |
| 2     | `TokenInterceptor` | Attaches `Authorization: Bearer <token>` when a token exists |

Angular runs them in registration order, so `ErrorInterceptor` is the outer layer: every response, successful or failed, reaches it before a caller's own error handler.

## What this means for code that calls HTTP

`ErrorInterceptor` handles a 401 by calling `UserService.logout()` synchronously — clearing the token, emitting a null user, and navigating to the login page — then rethrowing the error. By the time a caller's `error` handler runs, the session has already ended.

Consequences:

- A service's 401 handler cannot rely on the session still being alive. Treat 401 as "the session is already over" and 403 as "this call is not allowed for this user".
- Do not perform global auth transitions in feature services; that belongs to the interceptor and `UserService`. Feature code only reacts.
- The interceptor skips logout when the current route is `/register`, so registration failures do not bounce the user away.

`TokenInterceptor` calls `getToken()`, which fires the empty-token callback when no token exists; that callback also logs out, so a request made without a token ends the session too.
