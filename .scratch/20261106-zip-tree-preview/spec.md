# Zip Tree Preview

Status: ready-for-agent

## Problem Statement

A user looking at a Zip file on the manage-files page cannot see what is inside it without downloading it or unpacking it. Clicking the row does nothing, because a Zip file is not one of the previewable types. Downloading costs a round trip and leaves the archive on the user's machine. Unpacking writes thousands of rows and changes the directory, which is too much for the question "what is in this zip?". The user needs to inspect the archive's folder and file structure in place.

## Solution

A Zip file's expanded row gains a button that opens a dialog containing the archive's internal tree. The tree is read from the Zip central directory without decompressing any entry. Folders show their children. Files show their name and uncompressed size. The view is read-only. It writes no rows and extracts no files.

The tree is capped at 32 levels and 10,000 nodes. When a cap cuts the tree, the dialog says the tree was truncated and why. The response is cached for five minutes, so reopening the same archive is immediate.

Only the file's uploader can open the preview.

## User Stories

1. As a moon user, I want to open a Zip file's internal tree from the manage-files page, so that I can see what is inside without downloading it.
2. As a moon user, I want the tree to show folders and files with their names, so that I can recognize the archive's structure.
3. As a moon user, I want each file node to show its uncompressed size, so that I can tell a small config from a large binary.
4. As a moon user, I want folders ordered before files and each group sorted by name, so that I can find an entry predictably.
5. As a moon user, I want nested entries shown at the depth their path describes, so that the tree matches the archive.
6. As a moon user, I want explicit empty folders to appear, so that the tree shows the archive's full structure.
7. As a moon user, I want a Zip packed on Windows, with backslash separators, to render as a normal tree, so that the preview does not depend on the packing platform.
8. As a moon user, I want a deeply nested archive to show its first 32 levels and tell me the view was cut, so that the preview still opens.
9. As a moon user, I want an archive with more than 10,000 nodes to show the first 10,000 and tell me the view was cut, so that the dialog stays responsive.
10. As a moon user, I want to know when the tree is truncated, so that I do not mistake a partial tree for the whole archive.
11. As a moon user, I want the preview to open in a dialog from the expanded row, so that I keep my place in the file list.
12. As a moon user, I want the button to appear only on Zip files, so that I am not offered a preview that cannot work.
13. As a moon user, I want a corrupt or unreadable Zip to produce a clear error, so that I know the archive is the problem.
14. As a moon user, I want an empty Zip to show an empty tree, so that I can tell it contains no entries.
15. As a moon user, I want reopening the same Zip within five minutes to be fast, so that comparing entries does not cost a second read.
16. As a moon user, I want to see duplicate file entries when the archive contains them, so that the preview reflects the archive rather than a normalized view.
17. As a moon user, I want a file and a folder that share a path to both appear, so that neither is hidden.
18. As a moon user, I want archive names displayed as text, so that a crafted name cannot inject markup into the page.
19. As the file's uploader, I want only I can open the preview, so that archive names are not exposed to users who were granted access to download the file.
20. As an operator, I want the depth and entry caps to be configuration values, so that a deployment can tune them without a build.
21. As an operator, I want the preview to write nothing to the database or storage, so that a read never changes state.
22. As an operator, I want a cache failure to fall back to a direct read, so that Redis being down does not break the preview.
23. As an operator, I want the parser to read the central directory without decompressing entries, so that the preview costs milliseconds even on a large archive.
24. As a developer, I want the tree assembly to be a pure function, so that its edge cases are tested without a zip file.

## Implementation Decisions

**One read, two services.** vfm owns the user-facing preview endpoint, the owner check, and the cache. mini-fstore owns the archive bytes, the Zip parsing, the tree assembly, and the caps. vfm calls mini-fstore over its existing internal client and passes the caps it read from its own configuration.

**Read the central directory.** mini-fstore opens the stored file and reads its Zip central directory. It does not decompress entries, create temporary files, or write anything. A file that is not a readable Zip is refused with the existing illegal-format failure.

**Tree assembly sits in mini-fstore.** mini-fstore turns the flat entry list into one nested tree object and returns it to vfm. vfm does not rebuild the tree. It caches the object and returns it unchanged. This keeps archive parsing and the shape that depends on it in one service.

**Root and node shape.** The response carries one root folder node with an empty name. Its children are the archive's top-level entries. A folder node has a name, a folder flag, and children. A file node has a name, a file flag, and an uncompressed size. Folder nodes carry no size.

**Ordering.** Siblings are ordered folders first, then files, each group by name, case-insensitive. Two nodes with the same name keep their archive order relative to each other.

**Separators and names.** Both `/` and `\` are path separators. Leading separators and empty segments are ignored. Every other character is kept and displayed as text, including `..`, drive prefixes, control characters, and names longer than 255 bytes. Nothing is written, so no name can escape anywhere.

**Directory entries and duplicates.** An explicit directory entry becomes a folder node, so an empty folder survives. A folder inferred from a file path merges into the same folder node. A repeated directory entry merges into its folder. Every file entry becomes its own node, so a repeated file path renders twice. A file entry and a directory entry at the same path render as two nodes, one of each kind.

**Depth cap.** Depth counts path segments from the archive root, so `a/b/c.txt` sits at depth 3. Entries deeper than the configured maximum are left out and the response is marked truncated by depth.

**Entry cap.** The cap counts every node created in the tree, and the synthetic root is not counted. That means file nodes, explicit folder nodes, and folders inferred from file paths. A merged duplicate directory entry creates one node and counts once. A repeated file entry creates one node each time. The scan stops at the first entry that would create a node beyond the cap; entries that create no node, such as a repeated folder or an entry deeper than the depth cap, do not stop it. Later entries are left out, the response is marked truncated by entry count, and the total is not reported. The response reports how many nodes it included.

**Five-minute cache.** vfm caches the successful tree in Redis for five minutes, keyed by the stored file id. The cache read happens after the vfm row exists and belongs to the requester. Cached trees are returned unchanged, caps included. A failed preview is never cached. A cache read or write failure falls back to a direct mini-fstore call.

**Owner only.** The preview requires the requester to be the file's uploader. This is narrower than the download token rule, which also admits folder-granted users. The choice is deliberate. Archive entry names are metadata the uploader owns, and the preview exposes no bytes.

**No invalidation.** Stored file content is immutable for a given file id. Deleting or renaming the vfm row does not change the tree, and a deleted row fails the owner check before the cache is consulted. The five-minute expiry is the only bound on staleness.

**Frontend entry point.** The expanded row of a Zip file gains a Show tree button. The button fetches the tree and opens a dialog on success. On failure, the dialog does not open and the error message appears in the existing snackbar. The dialog renders the tree with the app's existing tree control, shows file sizes with the existing size formatter, displays the number of nodes the response included, and shows the truncation notice when the response is truncated. Top-level nodes start expanded. Deeper folders start collapsed.

**Zip detection.** The frontend offers the button for files whose name ends in `.zip`; a folder whose name ends in `.zip` is excluded. The backend does not trust the suffix. mini-fstore attempts to read the archive and returns the illegal-format failure when it cannot.

## Contracts

### Data

No schema change. The preview writes no rows and adds no columns.

The tree payload is JSON inside the response's data field. Node shape:

| Field | Type | Nullable | Meaning |
|---|---|---|---|
| `name` | string | no | One path segment, as stored in the archive. Empty for the root node only. |
| `dir` | boolean | no | True for a folder node, false for a file node. |
| `size` | integer | yes for folders and zero-byte files | Uncompressed size in bytes, from the central directory. Absent on folder nodes and zero-byte files. |
| `children` | array of nodes | yes for files and empty folders | Ordered children. Absent on file nodes and empty folders. |

Response data shape:

| Field | Type | Nullable | Meaning |
|---|---|---|---|
| `tree` | node | no | Root folder node. |
| `truncated` | boolean | no | True when a cap left nodes out. |
| `truncatedReasons` | array of strings | no | Zero or more of `depth` and `entries`. Empty when not truncated. |
| `includedEntries` | integer | no | Number of tree nodes represented, root excluded. |

### Interfaces

vfm endpoint, called by the frontend:

| Item | Value |
|---|---|
| Method and path | `GET /open/api/file/zip/tree` |
| Query field | `fileKey`, string, required. The vfm file uuid. |
| Authentication | Logged-in user, existing session. |
| Success body | Standard envelope with the response data shape above. |
| Failure `FILE_NOT_FOUND` | No live row for `fileKey`. |
| Failure `FILE_DELETED` | The row is logically deleted. |
| Failure `NOT_PERMITTED` | The requester is not the uploader, or `fileKey` names a folder. |
| Failure `ILLEGAL_FORMAT` | mini-fstore cannot read the file as a Zip. |

mini-fstore endpoint, called by vfm:

| Item | Value |
|---|---|
| Method and path | `POST /file/zip/tree` |
| Request `fileId` | string, required. The stored file id. |
| Request `maxDepth` | integer, zero means the server default. |
| Request `maxEntries` | integer, zero means the server default. Maximum tree nodes represented. |
| Success body | The response data shape above. |
| Failure `FILE_NOT_FOUND` | No stored file for `fileId`, or it is logically deleted. |
| Failure `ILLEGAL_FORMAT` | The file cannot be read as a Zip. |

mini-fstore applies no user authorization. vfm is the only caller and performs the owner check before this call, matching the existing internal boundary.

### Config

| Key | Type | Default | Owner | Effect |
|---|---|---|---|---|
| `vfm.zip-tree.max-depth` | integer | 32 | vfm | Deepest path segment count included in a preview. Sent to mini-fstore with each request. |
| `vfm.zip-tree.max-entries` | integer | 10000 | vfm | Maximum tree nodes represented in a preview. Sent to mini-fstore with each request. |

mini-fstore falls back to 32 and 10000 when a request field is zero or absent, and owns no preview configuration. A changed value applies to the next cache miss. A cached tree keeps the caps it was built with until it expires.

## Invariants

A preview writes nothing to the database, to storage, or to any other service.

A successful preview is returned only to the file's uploader.

A tree never contains a node whose path exceeds the depth cap, and never represents more nodes than the entry cap.

When any node is left out, `truncated` is true and `truncatedReasons` names at least one cause.

`includedEntries` equals the number of tree nodes in the returned tree, root excluded.

File sizes are uncompressed sizes from the central directory.

Sibling order is folders first, then files, each group case-insensitive by name. Duplicate names keep their archive order.

No preview request decompresses an entry or creates a temporary file.

## States and failure

A preview attempt ends in one of five outcomes.

**Served from cache.** The row exists, belongs to the requester, and a cached tree is present. The cached tree is returned without calling mini-fstore.

**Built and cached.** The row checks pass and no cached tree exists. mini-fstore reads the archive, builds the tree, and vfm caches the result before returning it.

**Truncated.** The archive is readable but a cap leaves nodes out. The response is a success with `truncated` true, and it is cached like any other success.

**Refused as unreadable.** mini-fstore cannot read the central directory. vfm returns `ILLEGAL_FORMAT`, and nothing is cached.

**Refused by access.** The row is missing, deleted, not owned by the requester, or names a folder. vfm returns the matching failure before reading the cache or calling mini-fstore.

A mini-fstore failure or timeout leaves nothing cached and surfaces as an error to the frontend. A Redis failure during cache read or write degrades to a direct call and does not fail the preview.

Deletion of the vfm row or the stored file while a tree is cached makes the next preview fail at the access check. The cached value expires unread.

## Boundaries

vfm owns `file_info`, the uploader check, the preview endpoint, and the cache. It reads the requester's identity from the existing session.

mini-fstore owns the `file` table and the stored bytes. It is authoritative for entry names, entry kinds, uncompressed sizes, whether the archive opens, and the tree shape. It holds no owner information and performs no authorization.

vfm passes both caps on every request, so mini-fstore holds no preview policy. A missing or non-positive value falls back to the same defaults, so the two services agree on the built-in values.

The cache key is the stored file id. The same content reached through two vfm rows shares one cached tree.

Only Redis is optional at preview time. Postgres, mini-fstore, and the vfm row are required. When mini-fstore is down the preview fails. There is no local fallback because vfm does not read archive bytes.

No other service is written to or read from. user-vault and the notification path are not involved.

## Testing Decisions

Good tests assert what a user or an operator can observe. That means the nodes and their order, the sizes, the truncation flags, the error code for an unreadable archive, and whether the cache or the archive read was used. They do not assert map iteration, internal struct names, or the order of writes that do not affect the result.

Four seams carry the feature.

**Tree assembly.** A pure function in mini-fstore takes a flat entry list plus the two caps and returns the tree, the truncation flags, and the included count. Table tests cover nesting, both separators, leading separators, `..` and control characters, explicit directory entries, empty folders, repeated file entries, repeated directory entries, a file and folder at one path, ordering, the depth boundary, the node boundary with inferred folders, and an empty archive. Prior art is the unpack plan tests.

**Archive reading.** One test builds a real Zip in a temporary directory with `archive/zip` and reads it through the new reader. It covers a normal tree, an explicit directory entry, a backslash path, a non-Zip file, a truncated file, and a cap small enough to truncate. Prior art is the unpack depth and temp directory tests.

**vfm handler.** Tests replace the mini-fstore client and the cache with fakes. They cover the owner path, a non-owner, a missing row, a deleted row, a folder key, a cache hit that skips mini-fstore, a cache miss that passes the configured caps, a cache write failure that still returns the tree, a cache read failure that still returns the tree, and an illegal-format failure that is not cached.

**Frontend.** The dialog component test renders the tree with a fixed response and asserts the node names, the sizes, the entry count, the truncation notice, the empty state, and crafted-name handling. The existing manage-files spec asserts that the Show tree button appears only for `.zip` rows, requests that row's file key, opens the dialog, and shows the snackbar on a refused or failed request. Prior art is the directory tree navigation spec and the manage-files spec.

No full-stack integration test is needed. The archive reader test covers the real parser, and the handler tests cover the service boundary with fakes.

## Out of Scope

Archives other than Zip, including tar, gz, 7z, and rar.

Reading inside a nested Zip.

Selecting entries, downloading a subset, or unpacking from the tree.

Reading or previewing an entry's contents.

Lazy loading of deeper levels. The whole tree within the caps arrives in one response.

Searching or filtering inside the tree.

Cache invalidation beyond the five-minute expiry.

Changing unpack behavior or any existing action in the expanded row.

Entry-count or uncompressed-size caps for unpack.

## Further Notes

The depth cap truncates rather than refusing, unlike unpack. A preview that shows nothing is less useful than one that shows the first 32 levels and says so.

Showing duplicate file entries keeps the preview faithful to the archive. Merging them would hide that extraction will overwrite one with the other. Folder entries merge because a folder is a container, and two container nodes with the same path would only split their children.

A crafted entry name is safe because the frontend renders names as text with Angular interpolation. The backend never turns a name into a path.

The cached tree is at most 10,000 nodes besides the root, roughly one to two megabytes of JSON, stored as one Redis string with no compression. The entry cap counts every node, including inferred folders, and is what bounds it. A cache fill reads the cache and writes the result with separate calls. The write takes the cache library's lock, but two concurrent first opens can both read the archive. That trade keeps a Redis failure from failing the preview, and the node cap bounds each read.

The related specs are `.scratch/zip-unpack-tree/spec.md`, which covers unpacking into real folder rows, and `.scratch/dir-zip-download/spec.md`, which covers the opposite archive direction. Neither is changed by this feature.
