# Dir Thumbnail Visibility

Status: ready-for-agent

## Problem Statement

A directory's cover is chosen by `dir-thumbnail` and `dir-thumbnail/batch` from the directory's
children, oldest first, requiring only `is_del = 0` and a non-empty `thumbnail`. The listing itself
hides children that are logically deleted or marked hidden, so the cover query can pick a child the
user cannot see. When that child sits in the trash, its thumbnail object has already been removed
from mini-fstore, the temporary-token call for it fails, and the directory reports no cover at all —
even though it holds live images that would serve fine.

Observed: directory `ZZZ2378659328458752615454` holds 1259 children; the oldest (id 81018) is
`is_logic_deleted = 1`, its thumbnail token minted nothing, and the batch response carried
`fstoreToken: ""` while live children from id 89771 onward had thumbnails.

## Solution

The cover query uses the same visibility rule as the listing: a child supplies a cover only when it is
not logically deleted, not physically deleted, and not hidden. The directory row itself must be
live as well. A directory whose only children are invisible reports no cover, which is correct — the
user sees an empty directory.

## User Stories

1. As a moon user, I want a directory's cover to come from a file I can actually see inside it, so
   that the preview never shows something the listing hides.
2. As a moon user, I want a directory whose visible children have no thumbnails to show no cover, so
   that the folder reads as empty rather than broken.
3. As a moon user, I want a directory of live images to get a cover again after trashed children are
   present, so that cleaning up a folder does not blank out its preview.
4. As a moon user, I want the hover tooltip and the browse-history cover to agree with the preview
   column, so that the same directory looks the same everywhere.

## Implementation Decisions

- The child predicate in both cover paths — the single-directory lookup and the batch lookup — gains
  the listing's visibility conditions: `is_logic_deleted = 0`, `is_del = 0`, `hidden = 0`, and a
  non-empty `thumbnail`. The pick order stays oldest `id` first.
- The directory row in the batch lookup also requires `is_logic_deleted = 0`; it already requires
  `is_del = 0` and the requesting user as uploader.
- No change to the response shape, the token-minting path, the ownership check, or the cover rule's
  ordering. A directory with no qualifying child keeps answering with an empty token, exactly as
  today.
- No schema, endpoint or config change.

## Testing Decisions

Good tests assert what the caller receives: which file id becomes the cover for a given set of rows,
and that the response is empty when only invisible children exist. They do not assert query text.

The existing backend tests for the directory thumbnail paths are the prior art. New cases: a child
that is logically deleted is skipped for a live one; a hidden child is skipped; a physically deleted
child is skipped; a child without a thumbnail is skipped; a directory whose children are all
invisible answers empty; the oldest qualifying child wins; the batch answers per directory without
letting one directory's trash change another's cover.

## Out of Scope

- Removing thumbnails of trashed files, or rewriting `thumbnail` when a file is deleted.
- Any change to the frontend preview column, which already falls back to the empty folder when a
  token is missing.
- Backfilling covers for directories that have none.
