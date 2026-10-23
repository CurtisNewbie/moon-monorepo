# 01: mini-fstore: bundle key and flat zip stream

**What to build:** A caller inside the platform can hand mini-fstore a directory name plus an ordered list of file ids and entry names, and receive a bundle key. Fetching the public archive route with that key downloads a zip holding exactly those entries, flat, with no directory entries inside it. The key expires after thirty minutes and refreshes while it is being read. Every id is resolved before the response body starts: an entry that cannot be served is skipped and logged, and the rest still download. An empty list yields a valid empty zip. The archive is named after the directory and entries are stored, not deflated.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] Minting a bundle key for an ordered list of existing files returns a key that resolves to the same entry names, in the same order.
- [x] The archive route streams a zip whose entries sit at the root: no directory entries, no nesting.
- [x] Each entry's bytes equal the stored file's bytes.
- [x] An empty file list produces a readable zip with zero entries.
- [x] A file whose record is absent, whose status is not normal, or whose stored bytes are missing is skipped and logged, and the remaining entries still download in order.
- [x] The response sets the archive content type and an attachment disposition carrying the directory name plus `.zip`.
- [x] The response carries no content length; the transfer is chunked.
- [x] Entries are stored rather than deflated.
- [x] A bundle key stops resolving after thirty minutes of no use, and its lifetime refreshes on a successful read.
- [x] A file entry's name is used verbatim; the id set and names come entirely from the bundle key.
