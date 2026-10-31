# 02: Folder preview toggle and cover-in-folder cell

**What to build:** A "Folder preview" checkbox in the manage-files toolbar, off by default and remembered in the browser. With it on, a directory row whose cover is already known draws the approved folder-with-cover glyph; with it off, or when no cover is known, the row keeps the plain folder glyph. No new request is made: the cover comes from the batch that already runs after each listing.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] The checkbox sits in the toolbar beside the existing thumbnail-tooltip checkbox, and is shown in the compact layout as well as the wide one
- [ ] A fresh browser shows it off, and every directory row draws the plain folder glyph
- [ ] Ticking it is remembered across reloads in that browser; an absent or unreadable stored value resolves to off
- [ ] With it on, a directory row whose cover is known draws the approved folder-with-cover glyph
- [ ] With it on, a directory whose cover is absent or has no image token keeps the plain folder glyph — no empty frame, grey box, or broken image
- [ ] With it off, a directory keeps the plain folder glyph even when its cover is already cached
- [ ] File rows render identically under both states
- [ ] Toggling issues no request, and listing, ordering, paging, actions, and the hover tooltip behave as before
- [ ] The label is translated in every locale the page ships
- [ ] Component tests cover the default, the persistence round-trip, cover present, cover absent, toggle off with a cached cover, and unchanged file rows
- [ ] Full frontend test suite and lint pass
