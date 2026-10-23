# Size Sort

Status: ready-for-agent

## Problem Statement

A user browsing a directory can order it by time, name, or a custom drag order. Nothing orders it by file size, so finding the files that take the most space means paging through the directory and reading the size column by eye.

Size has no meaning for a directory row, though. Ordering a mixed list by size would rank directories at zero bytes, so a size order only makes sense when the listing is files.

## Solution

Add `size` as a fourth sort field on the directory listing. It is offered only while the File type filter is selected, and choosing it makes the listing files only: a request that asks for size under any other type filter is coerced to File.

The listing opens largest first. The direction control still works. Size is transient, so the directory does not remember it and no Resume Cursor is recorded or offered under it.

The position endpoint resolves size as well, so jumping to a file from a search result or a deep link lands on the page the size-sorted listing actually shows.

## User Stories

1. As a moon user, I want to order a directory's files by size, so that I can find what takes the most space without paging through everything.
2. As a moon user, I want largest first by default, so that the space hogs are on page one.
3. As a moon user, I want to flip the direction, so that I can start from the smallest file when I am clearing out junk.
4. As a moon user, I want the size option only when I have filtered to File, so that I am never offered an order that cannot apply.
5. As a moon user, I want choosing size to show files only, so that the list matches the order I picked instead of mixing in directories ranked at zero bytes.
6. As a moon user, I want to switch the type filter back to All or Directory, so that I can browse directories again, and I accept that the sort returns to the directory's usual order.
7. As a moon user, I want equal-sized files to keep a stable order across pages, so that paging does not repeat or skip a file.
8. As a moon user, I want size order to work together with a name keyword search, so that I can rank only the matches.
9. As a moon user, I want a comic directory to keep opening in name-ascending order, so that reading order is unchanged.
10. As a moon user, I want the directory not to remember a size sort, so that my saved order for that directory is what greets me next time.
11. As a moon user, I want no resume offer under a size sort, so that I am not asked to continue at a position that only made sense for one size ranking.
12. As a moon user, I want clicking a search result while sorted by size to land on the page that holds that file, so that I do not have to hunt for it.
13. As a moon user, I want a deep link into a size-sorted directory to open the right page, so that bookmarks and shared links behave.
14. As a moon user, I want size order to respect the same visibility rules as every other order, so that hidden and deleted files never appear.
15. As a moon user, I want the file tray and bulk actions to work unchanged under size order, so that selecting across pages still works.
16. As a release engineer, I want this shipped in v0.0.6 without a schema change, so that no migration step is needed.

## Implementation Decisions

**Sort vocabulary.** `size` joins `time`, `name`, and `custom` as a sort field. It is a normal value of the existing order fields on both the listing and the position request, so no new endpoint or parameter is introduced. The deprecated `orderByName` flag is untouched.

**Availability and coercion.** Size applies only to a listing of files. When the listing request asks for size while the type filter is anything other than File, including absent, the effective type filter becomes File for that request. The response then contains files only, and the position query applies the same coercion. An unknown sort value keeps today's behavior and falls back to the default order.

**Ordering.** Size orders by `size_in_bytes` in the requested direction, then by row id in the same direction. Id breaks ties so pagination stays stable when several files share a size. There is no file-type prefix on the order clause, because the listing is files only.

**Default direction.** An omitted direction for size resolves to descending. An explicit ascending direction is honored.

**Comic directories.** A comic directory is forced to name-ascending for `time`, `name`, and `size`. A custom order is still honored, as today. The sort controls stay hidden there.

**Keyword search.** The name keyword filter keeps its current meaning and scope. Keyword filters the rows, size ranks them, and no relevance score is blended in. This matches how time with a keyword already ignores relevance.

**Transience.** Size is never persisted as an Order-By Preference. The client does not upload it on that path, so the stored preference for a directory stays whatever it was. Backend save behavior is unchanged.

**No Resume Cursor.** A listing sorted by size records no cursor, and no resume offer is shown for it. Existing cursors under time, name, or custom are untouched. The cursor's sort column never holds `size`.

**Position endpoint.** `file/position` understands size with the same ordering and the same File-type coercion. It serves file deep links and search-result jumps, not only cursor resume, so this is required even though size produces no cursor.

**Frontend behavior.** The sort dropdown shows the size option only while the File type filter is selected. Selecting size while the filter is All or Directory switches the filter to File. Changing the filter away from File resets the sort to the directory's stored or default order. Sort state continues to come from the server, and the results table is never re-sorted in the browser.

**Documentation and release.** The vfm API doc's sort enumeration notes the size value and its File-only rule. The change is recorded in the v0.0.6 release notes. No schema snapshot or change log entry is needed.

## Testing Decisions

Good tests assert visible behavior: the order of returned rows, which rows are returned at all, the page number a file lands on, and the request the client sends. They do not assert internal call order or private helpers.

Two existing seams cover the feature, and no new seam is needed.

The backend listing seam is the list operation itself, which already resolves sort and comic rules against a real test database. The existing sort tests are the prior art. New tests cover size ordering in both directions, the id tie-break, coercion of a non-File type filter to File, size under a keyword search, the comic force, and that a size listing returns no directories.

The backend position seam is the position calculation, tested the same way as the existing position tests. New tests cover the page a file lands on under size in both directions and under the coerced File filter.

The frontend seam is the manage-files component with its HTTP client stubbed. New tests assert that a size sort request carries the size sort with a File type filter, that choosing size under All or Directory sends File, and that leaving File resets the sort to the stored order. Existing component tests are the prior art.

## Out of Scope

- Sorting directory rows by size, or sorting a mixed list.
- Persisting size as an Order-By Preference or resuming under it.
- Schema or index changes.
- Showing human-readable or relative sizes differently from today.
- A global largest-files view across directories.
- Changing how the custom drag order or comic reading order works.

## Further Notes

Known limitation. A size sort is a filesort on the server, both inside one directory and over a keyword-matched set. No index helps the keyword path, and a size index buys nothing for the in-directory case, so this ships without one. Revisit if a real directory or search proves slow.

Target release is v0.0.6.
