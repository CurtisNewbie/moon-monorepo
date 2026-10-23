# Directory Zip Download

Status: ready-for-agent

## Problem Statement

A user who wants a whole directory today has to open it and download each file one at a time. vfm refuses a directory outright: the download token endpoint rejects anything that is not a file ("Downloading a directory is not supported"). The only archive operation in the platform runs the other way, unzipping an uploaded archive into a directory.

## Solution

A directory's expanded panel gets a zip download action, shown only for rows that are directories. Clicking it hands the browser a link that streams a zip built on the fly from the directory's own files.

The zip is flat: it holds the directory's direct file children, each at the archive root, and nothing else. Sub-directories are skipped, and their contents never enter the archive. A directory with no direct files still produces a valid empty zip.

Only the directory's owner can download it.

## User Stories

1. As a moon user, I want to download a directory as one zip file, so that I do not fetch its files one by one.
2. As a moon user, I want the action offered only on directory rows, so that file rows keep their existing single-file download.
3. As a moon user, I want the archive to contain exactly the files I see directly in that directory, so that what I get matches what I browsed.
4. As a moon user, I want sub-directories left out of the archive, so that the download stays flat and predictable.
5. As a moon user, I want the archive named after the directory, so that I can tell two downloads apart in my downloads folder.
6. As a moon user, I want two files that share a name inside the flat archive to both survive, so that neither is silently lost.
7. As a moon user, I want an empty directory to download as a valid empty zip, so that the action behaves the same everywhere.
8. As a moon user, I want the download to start immediately, so that I do not wait for an archive to be built and stored first.
9. As a moon user downloading a large directory, I want the transfer to stream, so that memory use does not grow with the directory's size.
10. As a moon user, I want the download to reuse the browser's normal download flow, so that I get the progress and save dialog I already expect.
11. As a moon user with Chinese characters in my file and directory names, I want them preserved in the archive and its filename, so that nothing is mangled.
12. As a moon user, I want files that are already compressed to be stored without recompression, so that a download is not slowed for no gain.
13. As a moon user, I want a directory I deleted to stop downloading, so that the action follows the same rules as the rest of the app.
14. As a moon user without ownership of a directory, I want the download refused, so that no one else's files are exposed to me.
15. As a moon user whose directory holds a file whose stored bytes have gone missing, I want a clear failure instead of a half-written archive, so that I never keep a silently truncated zip.
16. As a moon user, I want the download link to stop working after half an hour, so that a copied link does not stay valid forever.
17. As a moon user, I want the tooltip on the action to say what it does, so that I do not confuse it with the unzip action.
18. As an operator, I want the zip built only from validated file records, so that a deleted or invalid record never reaches an archive.

## Implementation Decisions

**Where the archive is built.** In mini-fstore, next to the bytes. Its file transfer routine already writes to any `io.Writer`, so a zip writer can consume it directly. No archive is persisted: the response streams, and nothing is written to storage or to a temp directory.

**Selection rule.** vfm reads the directory's direct children: same parent, file type file, not deleted. Nothing recurses. The resolved set is ordered by name ascending, so the same directory always produces the same archive order.

**Permission.** The directory must exist, must be a directory, must not be logically deleted, and its uploader must be the caller. Virtual folder grants are not consulted, because a grant only ever covers file rows and never a directory row.

**Token handoff.** The existing temp token store holds exactly one file, so it cannot express a directory. A parallel bundle entry is added to the same Redis store under its own prefix. It carries the archive filename (the directory name) and an ordered list of file id plus entry name. Lifetime is thirty minutes and it refreshes on each read, matching the single-file key. The bundle key is a bearer secret, the same way a single-file key is today.

**Endpoints.** vfm gains one authenticated action that takes a directory file key and returns the bundle key. mini-fstore gains one public route that takes the bundle key and streams the archive. The existing single-file token endpoint keeps rejecting directories.

**Validation before the first byte.** The whole file set is checked in one batch before the response body starts. If any record is missing or not normal, the request fails with an error status and no zip bytes. This is what keeps a partial archive from ever reaching the user, since a stream cannot take back what it already sent.

**Response shape.** Also before the body starts, the response sets the zip content type and an attachment disposition carrying the directory name plus `.zip`. Content length is deliberately absent, because the archive size is unknown until the last entry is written; the transfer is chunked. Entry names are bare filenames, so the archive carries no directory entries at all.

**Entry naming.** A name that already appears earlier in the ordered set gets ` (n)` appended before its extension, counting from 1. Because the order is deterministic, the same collision always resolves the same way.

**Compression.** Entries are stored, not deflated. Directory contents here are largely images and other already-compressed media, and the host runs one CPU with a 300 megabyte memory cap.

**Frontend.** The action lives in the expanded panel's directory action block, next to the rename, comic toggle, and truncate actions, as an icon button with a tooltip read from the translation catalog. The handler asks vfm for the bundle key, then opens the fstore archive route the same way the single-file download opens its own route. No new component or dialog.

**No schema change.** Nothing is persisted, so no table, snapshot, or release increment changes.

## Testing Decisions

Tests assert external behavior only: what a caller receives, and what ends up inside the archive. No test asserts on the internal shape of the Redis entry.

Three seams, all existing:

- mini-fstore's storage package, driven through an `httptest` request against a temporary storage directory, the pattern its range-serving test already uses. Covers: entries appear flat at the archive root, store method, byte content matching the source files, a single entry, an empty set producing a readable empty zip, and a bad entry failing before any body byte.
- vfm's service layer against real MySQL, the pattern its existing service tests already use. Covers: direct file children selected and sub-directories excluded, name ordering, the ` (n)` collision suffix, a deleted directory refused, a non-owner refused, an empty directory yielding an empty bundle, and a bundle key that resolves back to the same names and ids.
- The manage-files component with its hand-rolled HTTP fake, the pattern its existing spec uses. Covers: the action posts the directory's file key and opens the returned archive route, and nothing else changes about the file-row download path.

Prior art: mini-fstore's serving test for the streaming seam, vfm's service tests for the selection and permission seam, the manage-files spec for the trigger seam.

## Out of Scope

- Recursive archives, or any option to include sub-directories.
- Downloading a directory for a virtual folder grantee who does not own it.
- Pre-building and caching an archive as a stored file, with progress reporting or a notification when it is ready.
- Downloading several selected rows, or a whole virtual folder, in one archive.
- Any change to the single-file download, its token, or the unzip action.
- Archive formats other than zip.

## Further Notes

- Backend work loads the `miso` skill first, per the backend guide.
- The tooltip goes through the app's translation catalog rather than a literal string.
- The frontend diff is verified as the frontend guide requires.
