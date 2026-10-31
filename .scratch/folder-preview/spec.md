# Folder Preview

Status: ready-for-agent

## Problem Statement

The manage-files list shows two kinds of row in one Preview column. A file row shows its thumbnail;
a directory row shows a plain folder glyph, so every directory looks identical no matter what is
inside it. The user cannot tell a directory full of screenshots from an empty one without opening
each in turn, and the folder glyph carries no information at all.

The directory's own cover image already exists on the server: the page fetches it in batch for the
hover tooltip, but the list never draws it. The information is on screen and thrown away.

## Solution

A "Folder preview" checkbox in the manage-files toolbar. While it is on, a directory row whose cover
image is known draws that image inside the folder glyph, so the row reads as a directory that has a
picture in it. While it is off — and for a directory with no cover image — the row keeps today's
plain folder glyph.

The checkbox starts off, so nobody's list changes until they ask for it, and the choice is remembered
in the browser.

## User Stories

1. As a moon user, I want an on/off switch for folder previews, so that I can choose between a
   compact familiar list and an informative one.
2. As a moon user, I want the switch to start off, so that my list looks the same after this ships as
   it did before.
3. As a moon user, I want my choice remembered across reloads, so that I do not re-tick it every
   visit.
4. As a moon user, I want a folder row to show a picture when its folder has one, so that I can
   recognise a directory by content rather than by name.
5. As a moon user, I want the picture drawn inside a folder shape, so that I never mistake a
   directory for an image file.
6. As a moon user, I want a folder with no picture inside to keep the plain folder glyph, so that
   empty and non-picture folders do not look broken.
7. As a moon user, I want the picture to appear as soon as the page knows it, so that I am not
   waiting on a second round-trip I already paid for.
8. As a moon user, I want the preview to use the same cover the tooltip and the browse-history page
   use, so that the same folder looks the same everywhere.
9. As a moon user on a phone, I want the same switch and the same previews, so that the compact list
   is not a degraded version.
10. As a moon user, I want file rows unchanged, so that a fix for folders does not disturb image
    thumbnails that already work.
11. As a moon user, I want the hover tooltip to keep behaving as before, so that the two switches
    stay independent and I can have either, both, or neither.
12. As a moon user, I want the switch to be one click away in the toolbar, so that toggling does not
    mean opening a settings page.
13. As a moon user with many folders on a page, I want no extra requests when I tick the box, so that
    turning previews on is instant rather than a page of new fetches.

## Implementation Decisions

**Frontend only.** The directory cover already exists and is already fetched: after every listing the
page asks the batch directory-thumbnail endpoint for every directory uuid on the page and caches the
resulting image URL per uuid. This feature changes what the Preview cell draws; it adds no endpoint,
no request, and no backend change.

**Cover rule is the existing one — reused, not redefined.** The cover is the oldest direct child of
the directory that has a generated thumbnail. Subdirectories are not descended into. Because
thumbnails are generated for video as well as image children, a video child can be a directory's
cover; that is accepted, not special-cased. The tooltip, this column, and the browse-history page all
read the same rule.

**Toggle.** One boolean, "Folder preview", default off, persisted in browser localStorage so it
survives reloads and applies per browser rather than per account. It is a new, separate control: the
existing "Enable Thumbnail Tooltips" checkbox keeps its current behavior and its current
component-only (unpersisted) state.

**Placement and reach.** The checkbox sits in the toolbar beside the tooltip checkbox. Unlike the
tooltip checkbox it is shown in the compact layout as well as the wide one, because the compact
layout keeps the Preview column.

**Cell rendering.** For a directory row while the toggle is on and a cover URL is known, the cell
draws the cover image clipped inside a folder-shaped frame whose tab stays visible: the row must read
as a directory holding a picture, not as a picture. For a directory row with no known cover, and for
every directory row while the toggle is off, the cell draws today's folder glyph. File rows are
untouched.

**Timing.** Covers arrive from the batch request that follows each listing, so with the toggle on a
row may draw its folder glyph first and swap to the composite once the batch answers. No placeholder,
spinner, or layout change is introduced for that gap; a directory whose batch entry is absent or has
an empty token keeps the glyph.

**Cache behavior.** The cover cache is keyed by directory uuid and already survives repeated listings
and paging, so revisiting a page draws covers immediately. No eviction strategy is added.

**Labels.** The checkbox label is translated in every locale the page ships.

## Testing Decisions

Good tests assert what the user sees and what the browser remembers: whether a directory row renders
an image or a glyph, and whether the choice persists. They do not assert CSS, DOM class names, or
which private method computed a URL.

One seam covers the feature: the manage-files component with its HTTP client stubbed, matching the
existing component tests. New tests cover:

- the toggle defaults to off on a fresh browser and the cell draws the folder glyph for directories;
- ticking it persists the choice, and a fresh component instance reads it back as on;
- with the toggle on and a cover present in the batch response, a directory row draws an image;
- with the toggle on and the batch returning no token for a directory, that row keeps the folder
  glyph;
- with the toggle off, a directory row draws the folder glyph even when a cover is cached;
- file rows render exactly as before under both toggle states.

No backend test is added: the cover rule and endpoint are unchanged, and the existing tests for them
stand.

## Out of Scope

- Changing the cover rule (recursive descent, explicit user-set cover, images-only children).
- Any backend, API, or schema change.
- Cover previews in the expanded row detail, the directory browse-history page, or the directory
  tree.
- Suppressing or changing the hover tooltip when previews are on.
- Persisting the setting per account or to the server; syncing it across browsers.
- Lazy loading, virtualization, or preloading of cover images.
- A placeholder or skeleton for directories whose cover has not arrived yet.

## Contracts

### Data

N/A — no stored data changes. The only new persisted value is a browser-local boolean preference
(see Config), which is not a server-side record and carries no migration.

### Interfaces

N/A — no new or changed API. The feature consumes two existing contracts as they are:

- the listing response already carries each directory's `uuid`, which the page batches;
- the batch directory-thumbnail request takes a list of directory keys and answers one entry per
  directory, each with the directory key and an image token that is empty when the directory has no
  cover. A directory with an empty token is indistinguishable, on the client, from a directory the
  batch did not answer for — both draw the glyph.

Failure of the batch request is silent today and stays silent: no cover is shown, the glyph stands,
and the row is otherwise unaffected. A missing cover is never an error state the user sees.

### Config

| Key | Type | Domain | Default | Owner | Effect |
|---|---|---|---|---|---|
| folder preview preference | boolean, localStorage | on / off | off | manage-files page | on: draw the cover inside the folder glyph for directory rows; off: today's folder glyph |

Absent or unreadable stored value resolves to off. The key is namespaced to the manage-files page so
no other page reads it.

## Invariants

- A directory row is always drawn as a folder. The cover, when present, is inside the folder shape,
  never a bare image.
- A directory with no cover is never drawn as an empty frame, a grey box, or a broken image.
- The preview toggle never changes which rows are listed, their order, their paging, or their
  actions.
- File rows are identical under both toggle states.
- Toggling the checkbox issues no request.

## States and failure

The switch has two states and no derived or disabled state. The cell has two states per directory
row: glyph, or glyph-with-cover. A row moves glyph → glyph-with-cover when its cover becomes known,
and back only when the toggle turns off or a later listing leaves it without a cover. The component
holds no "loading cover" state. Unmounting the page clears the cover cache; the preference outlives
it.

## Boundaries

N/A — the feature crosses no system boundary. The browser-local preference is read and written by the
page alone; the cover images and their tokens are served by existing endpoints under their existing
auth model.

## Further Notes

Target release: next frontend release; no schema snapshot or change-log entry is involved.
