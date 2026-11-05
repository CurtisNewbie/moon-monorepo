# Zip Unpack Tree

Status: ready-for-agent

## Problem Statement

Unpacking an archive loses its folder structure. Every entry lands as a direct child of the one directory vfm creates, and its name becomes the whole path it had inside the archive. A file that sat at `Documents/doc/jdk-12.0.2_doc-all/docs/api/jdk.xml.dom/org/w3c/dom/html/HTMLBRElement.html` shows up as a single file with that entire string as its name, one hundred and six characters long.

Three consequences follow. A user cannot browse what they uploaded; they scroll a flat list of long, near-identical names. A path longer than 255 characters cannot even be stored, because the name column holds 255. And the archive's own directory entries are turned into zero byte files, so `Documents/` and `Documents/doc/` appear as junk records next to the real content.

Separately, nothing stops an archive whose tree is absurdly deep. Such a tree is useless to browse, and building it costs rows and time across two services.

## Solution

Unpacking recreates each entry's folder layers under the target directory, so the archive's own structure becomes the directory structure the user browses.

The archive's directory entries are no longer stored. Only files create rows, and their parent directories are created on the way down. A folder that holds no files does not survive, because nothing creates it.

Before any extraction starts, the depth of the archive is checked. An archive whose deepest entry sits more than the configured number of levels down the tree is refused whole. No extraction happens, no rows are written, and the user is told why. An archive that cannot be read is refused the same way. In both cases the request fails with an error, the empty target directory stays, and one notification explains what happened.

Names inside an archive that cannot be honoured, such as a path escaping the target directory or a segment too long for the name column, are skipped one by one. The rest of the archive still unpacks, and one notification reports how many entries were left out.

The operation stays asynchronous. The caller learns of a refusal immediately, and learns of the finished tree through the files appearing in the directory.

## User Stories

1. As a moon user, I want the folders inside my archive to become folders in the target directory, so that I can browse the archive the way it was packed.
2. As a moon user, I want a deeply nested archive to keep its full depth, so that source trees and documentation sets arrive intact.
3. As a moon user, I want each file to keep only its own name, not its whole path, so that names stay readable and short.
4. As a moon user, I want files nested more than five levels deep to open normally, so that nothing in the archive becomes unreachable through the browser.
5. As a moon user, I want an archive's empty folders not to appear, so that the tree I get contains only what I can open.
6. As a moon user, I want garbage directory records from archive folder markers to stop appearing, so that the target directory holds my files and nothing else.
7. As a moon user, I want an archive that is too deeply nested to be refused before anything is written, so that I never end up with half a tree.
8. As a moon user, I want a clear notification when my archive is refused for depth, naming the limit and the deepest path, so that I know what to change.
9. As a moon user, I want a corrupt or truncated archive refused the same way, so that I am not left guessing why nothing appeared.
10. As a moon user, I want the refusal to also surface on the action itself, so that I see it without opening my notifications.
11. As a moon user, I want an entry with a path that escapes the target directory dropped, so that an archive cannot write outside the folder I chose.
12. As a moon user, I want an entry with an unusable name dropped rather than rewritten, so that the files I do get have the names the archive gave them.
13. As a moon user, I want one notification listing how many entries were dropped and which ones, so that I can spot what is missing.
14. As a moon user, I want the rest of the archive to unpack even when a few entries are dropped, so that one bad entry does not cost me the whole upload.
15. As a moon user, I want an entry whose folder already exists in the target directory to merge into it, so that unpacking twice into the same place does not build a second copy of the tree.
16. As a moon user, I want a file whose name is already taken to be stored under a numbered name, so that neither file is lost.
17. As a moon user, I want a blocking file in the target directory to make room for a folder through a numbered name instead of dropping the subtree, so that the archive's contents survive a name clash.
18. As a moon user, I want numbered names to follow the same convention as the zip download action, so that the two directions of archiving look alike.
19. As a moon user, I want the target directory named after the archive as it is today, so that nothing changes about where the result appears.
20. As a moon user, I want the unpacked files to belong to me and sit in the directory I chose, so that permissions match everything else I own.
21. As a moon user, I want a failure in the middle of an unpack to stop further work, so that a database problem does not produce thousands of failed rows.
22. As a moon user, I want to be told when an unpack stopped early and how many entries were not created, so that a partial tree is never mistaken for a complete one.
23. As a moon user, I want a partial tree left as it is, so that I can delete it myself once I have read the notification.
24. As a moon user, I want the notification to arrive once per unpack, so that a thousand dropped entries do not mean a thousand messages.
25. As an operator, I want each dropped entry logged with its path and reason, so that I can tell why a tree is incomplete.
26. As an operator, I want the depth limit to be a configuration value, so that a different deployment can tighten or loosen it without a build.
27. As an operator, I want a database failure during unpack logged at ERROR level and the message acknowledged, so that the queue is not republished forever.
28. As an operator, I want the archive inspected for depth without decompressing it, so that the check costs milliseconds on a large archive.

## Implementation Decisions

**One operation, two services.** vfm owns the directory tree and runs the user-facing unpack action. mini-fstore owns the archive bytes and performs the mechanical unzip. The term Unpack names the operation; unzip names the step inside mini-fstore. Nothing is renamed.

**Where the depth check runs.** In mini-fstore, inside the synchronous call vfm already makes to start an unpack. The check reads the archive's central directory, which needs no decompression. A refused archive returns an error from that call and nothing is enqueued, so no extraction ever starts. vfm turns that error into a notification and also returns it to the caller.

**The limit travels with the request.** vfm reads its own configuration and sends the maximum depth in the unpack request. mini-fstore applies the value it is given, and falls back to the same default when a caller omits it. mini-fstore holds no policy of its own.

**What depth counts.** The number of path segments in the deepest file entry, counted from the target directory, with the archive's own top level folder counting as level one. Directory entries are ignored, since they are never created. Empty segments are ignored, and an entry containing `..` counts as one segment more than it looks.

**Depth only.** No cap on entry count and no cap on total uncompressed size in this change.

**Tree rebuild.** vfm rebuilds the tree from the reply it already receives. Each entry name is sanitised into segments, ancestors are created once and reused across entries, and the file is saved under its deepest parent. Existing children of the target directory are read once before the loop so a folder that already exists is merged rather than duplicated.

**Reusing helpers for writes.** Rows are created through the existing directory and file helpers rather than a bulk insert. This keeps one code path for file creation and one for directory creation, at the cost of a lock and a cycle check per file. An archive with thousands of entries is correspondingly slow, and that cost is accepted here.

**Name rules.** A name is honoured only if it can be expressed. It carries no path separator inside a segment, no `..`, no control character, no trailing dot or space, and no segment longer than 255 bytes. An entry that fails is dropped, never rewritten, matching the rule the zip download action already follows. Path separators at the start are dropped, and backslashes are read as separators, so an archive packed on another platform still unpacks.

**Collisions.** A folder whose name already exists as a folder is reused. A folder whose name is held by a file gets the numbered suffix. Two entries resolving to the same file path get the numbered suffix on the later one, before the extension, through the same helper the zip download action uses. Numbering starts at one and depends on entry order.

**Directory entries.** mini-fstore stops creating stored files for entries that are directories, so they never reach the reply and never become rows.

**Failure inside the asynchronous step.** The first error while writing rows aborts the loop. The handler logs it at ERROR level and returns no error, so the message is acknowledged and never redelivered. The rows already created stay.

**Notifications.** vfm sends notifications through user-vault's create endpoint. One notification per unpack at most, combining the dropped entry count and the incomplete count when both apply. A notification that cannot be sent is logged and never fails the unpack.

**Kept as is.** The target directory keeps its current name, built from the archive name plus a timestamp. Ownership, deletion state and emptiness checks on the target directory are unchanged. The reply payload keeps its shape, because a synchronous error now carries every refusal.

## Contracts

### Data

No schema change. `file_info` rows keep their current shape.

Directory rows carry the directory file type, size zero, no fstore file id, the target directory or an intermediate folder as parent, and the caller as uploader. File rows carry the file file type, the fstore file id, the size reported by mini-fstore, and the deepest folder as parent.

The 255 byte limit on the name column governs each path segment, not the whole path.

### Interfaces

`POST /file/unzip` in mini-fstore gains one request field, the maximum depth, an integer where zero means the server default. Its response stays empty on success. It gains one failure code for a tree past the limit, and reuses the existing illegal format code for an archive that cannot be read. An unpack with a file that is not a zip keeps returning the not a zip failure as today.

The unpack reply keeps its shape. It carries the archive file id, the entry list with name, size, checksum and fstore file id, and the caller's extra payload.

The vfm unpack action keeps its request. It now fails with the refusal error when mini-fstore refuses the archive.

The notification call is user-vault's create endpoint, with a title up to 255 characters, a message up to 1000, and the receiving user number.

### Config

`vfm.unpack.max-depth`, an integer, default 32, owned by vfm. It is the effective limit for every unpack, and each change applies from the next unpack. mini-fstore has no new knob; its fallback for a missing value is also 32.

## Invariants

No row is ever created outside the target directory, whatever the archive contains.

Every row created by an unpack carries the requesting user as uploader, and the target directory or one of its descendants as parent.

No directory row exists for an archive directory entry, and no file row has a path separator in its name.

An unpack that is refused creates no rows at all.

An unpack never reports completion. The tree appearing in the directory is the only success signal.

At most one notification is sent per unpack attempt.

## States and failure

An unpack attempt is in one of four outcomes.

**Refused by depth.** Detected before enqueue. No rows are written, the target directory stays empty, the caller gets an error, and one notification is sent.

**Refused as unreadable.** Same shape as depth, with a different reason in the notification.

**Completed with drops.** Every honourable entry is written, dropped entries are counted per reason, and one notification is sent when the count is above zero.

**Stopped early.** A row write failed, the remaining entries were not attempted, the partial tree is kept, and one notification says how many entries were not created.

The refusal states are terminal. The stopped state is terminal too, because the message is acknowledged rather than retried; nothing repairs a partial tree except a fresh unpack by the user.

A notification failure is logged and does not change the outcome.

## Boundaries

vfm writes every vfm row and every notification. mini-fstore writes the archive's stored bytes and its own file rows.

mini-fstore is authoritative for anything read out of the archive, meaning entry names, sizes, checksums, and whether the archive opens at all. vfm trusts nothing about the archive it cannot see.

The synchronous unpack call is the only place requiring both sides to agree, and its contract is the request field plus the failure codes. When mini-fstore is down, the call fails like any other network failure and the unpack does not start.

The reply message is asynchronous and carries no retry header, so vfm's acknowledgement decision is final. Ordering between entries is the order mini-fstore packed them, and vfm applies numbering in that order.

user-vault is written to, never read, by this flow. Its being down costs the notification and nothing else.

## Testing Decisions

Good tests here assert what a user would observe, meaning which folders and files exist after an unpack, under which parent, and with which names. They do not assert the order in which rows were inserted beyond what naming depends on.

One seam carries the path logic. A pure planner takes the entry names from a reply plus the target directory's existing children, and returns the folders to create, the files to create with their resolved names, and the entries dropped with a reason. It is unit tested with no database, in the style of the zip download action's naming tests, which already cover the shared numbering helper.

A second pure seam covers depth. A function takes entry names and returns the deepest level, tested with table cases including directory entries, empty segments and paths escaping the root.

mini-fstore's archive reading is covered by one test that builds a real zip in a temporary directory, the way the existing bundle tests build theirs, and asserts the refusal for an archive past the limit and the pass for one at it.

The database writer stays thin and untested on its own, since it delegates to helpers that already have integration coverage. The existing unpack integration test keeps working and is not part of the default run.

## Out of Scope

Unzipping straight into storage without the temporary directory. The staging directory stays.

Caps on entry count or total uncompressed size.

Recreating empty folders from the archive.

Repairing trees produced by earlier unpacks, which stay flat.

Any frontend change, since the tree already renders as folders once the rows exist.

Retrying the asynchronous step, or making it idempotent.

## Further Notes

The depth limit exists to keep a browse-able tree, not to defend storage. A shallow archive with ten million entries still gets through, and the accepted write path makes it slow. A size and count guard is a separate decision.

Skipping entries rather than rewriting them is deliberate and matches the zip download direction. It means a Windows archive full of backslashes unpacks, while a name carrying a control character does not.

Because the asynchronous step is acknowledged on failure, a transient database outage costs the rest of an archive. The notification is what makes that visible.

Related, `.scratch/dir-zip-download/spec.md` covers the opposite direction and owns the numbering convention this spec reuses.
