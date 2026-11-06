# 02: Serve and cache the Zip tree preview in vfm

**What to build:** vfm gains the logged-in endpoint the manage-files page calls to preview a Zip file. The requester must be the file's uploader. A missing row returns `FILE_NOT_FOUND`, a logically deleted row returns `FILE_DELETED`, and a row that is not owned by the requester or names a folder returns `NOT_PERMITTED`.

The endpoint reads the depth and entry caps from vfm configuration and passes both to mini-fstore with the stored file id. A successful reply is cached in Redis for five minutes under the stored file id and returned unchanged, caps included. A cache hit skips the mini-fstore call. A refused preview is never cached. A cache read or write failure falls back to mini-fstore and does not fail the request.

mini-fstore's `FILE_NOT_FOUND` and `ILLEGAL_FORMAT` failures reach the caller as the same codes. No schema changes and no rows are written.

**Blocked by:** 01, Read a Zip's tree in mini-fstore

**Status:** ready-for-agent

- [ ] The endpoint requires a logged-in session and accepts a vfm file key.
- [ ] Only the file's uploader is served. A non-owner and a folder key return `NOT_PERMITTED`, a missing row returns `FILE_NOT_FOUND`, and a deleted row returns `FILE_DELETED`.
- [ ] `vfm.zip-tree.max-depth` defaults to 32 and `vfm.zip-tree.max-entries` defaults to 10000. Both are sent on every mini-fstore request.
- [ ] A successful reply is cached for five minutes under the stored file id. A cache hit returns the cached tree without calling mini-fstore.
- [ ] A failed preview is not cached.
- [ ] A cache read or write failure still returns a freshly built tree.
- [ ] mini-fstore's missing-file and illegal-format failures reach the caller unchanged.
- [ ] Handler tests with fakes cover owner, non-owner, missing row, deleted row, folder key, cache hit, cache miss with caps, cache read failure, cache write failure, and illegal format not cached.
