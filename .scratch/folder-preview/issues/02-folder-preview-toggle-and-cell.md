# 02: Folder preview toggle and cover-in-folder cell

**What to build:** A "Folder preview" checkbox in the manage-files toolbar, off by default and remembered in the browser. With it on, a directory row whose cover is already known draws the approved folder-with-cover glyph; with it off, or when no cover is known, the row keeps the plain folder glyph. No new request is made: the cover comes from the batch that already runs after each listing.

**Blocked by:** 01

**Status:** done

- [x] The checkbox sits in the toolbar beside the existing thumbnail-tooltip checkbox, and is shown in the compact layout as well as the wide one
- [x] A fresh browser shows it off, and every directory row draws the plain folder glyph
- [x] Ticking it is remembered across reloads in that browser; an absent or unreadable stored value resolves to off
- [x] With it on, a directory row whose cover is known draws the approved folder-with-cover glyph
- [x] With it on, a directory whose cover is absent or has no image token keeps the plain folder glyph — no empty frame, grey box, or broken image
- [x] With it off, a directory keeps the plain folder glyph even when its cover is already cached
- [x] File rows render identically under both states
- [x] Toggling issues no request, and listing, ordering, paging, actions, and the hover tooltip behave as before
- [x] The label is translated in every locale the page ships
- [x] Component tests cover the default, the persistence round-trip, cover present, cover absent, toggle off with a cached cover, and unchanged file rows
- [x] Full frontend test suite and lint pass

## Comments

Implementation notes:

- Tests: full suite green (`npm test -- --watch=false --browsers=ChromeHeadless` → 116 SUCCESS). `npm run lint:breakpoints` passes.
- `npm run lint` (breakpoints + `ng lint`) is red on the integration tip itself, before this change: the
  repo-wide TSLint profile demands single quotes and brace/semicolon rules the code base does not
  follow (328 pre-existing errors in `mng-files.component.ts` alone). No new violations were added:
  the touched TS file and spec hold their baseline error counts, `translate.dict.ts` gains only the
  two new key lines in that file's prevailing double-quote style. So the lint half of the last
  criterion is unmet for reasons outside this ticket; the test half passes.
- Closed 2026-10-31 with the user's agreement: the box is ticked on the test half, and the lint half
  stays red only on the repo-wide pre-existing TSLint rot (3393 errors on untouched files).
- Storage key: `manage-files.folderPreview`, values `'1'` / `'0'`. Read in `ngOnInit`, written from
  `(ngModelChange)`; a browser that throws on `localStorage` leaves the toggle off / session-only.
- Rendering is gated by `showDirCover(f)` = `f.isDir && folderPreview && hasThumbnail(f)`.
  The composite markup and CSS follow the approved recipe in ticket 01's comments; file rows keep
  their exact previous markup.

Review rework (2026-10-31):

- Expired/broken cover: the cover `<img>` reports `(error)` to `onDirCoverError(f)`, which deletes the
  uuid from `dirThumbnailCache` so the cell falls back to the plain folder glyph on that change
  detection pass (event bindings run in the zone) and on later listings. Removed
  `background: #2a2a2a` from `.dir-cover-img`: it only existed to backstop a missing picture, which
  now cannot persist, and while loading it would flash a dark box over the folder body. Test:
  `falls back to the plain glyph when a cached cover fails to load`.
- Compact reachability: the toggle now carries `.folder-preview-toggle`, and the existing
  `@media (max-width: 768px)` block gives it `order: -1`, so on a phone it is first in the
  horizontally scrolling search strip instead of off-screen at its end. Wide layout unchanged (still
  beside the tooltip checkbox). The compact test now sizes the Karma runner iframe to 412px and
  asserts the checkbox rect lies inside both the strip's visible area and the viewport; it fails when
  the `order` rule is removed.
- Test hygiene: the spec's `flushAsync()` helper moved to module scope and is reused; `beforeEach`
  removes only `manage-files.folderPreview`; `thumbnailCell()` lost its unused index parameter. The
  stub child components keep every member they declare because the real template binds all of them.
- `showDirCover` reuses the existing `hasThumbnail(f)` instead of re-reading the cache itself.

