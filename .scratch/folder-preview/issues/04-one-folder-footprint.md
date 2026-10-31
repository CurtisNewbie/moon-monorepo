# 04: One folder footprint, a legible cover

**What to build:** Every directory row draws the same tabbed folder shape at the height of a file row's thumbnail, so a column of directories and files lines up row for row and nothing jumps when a cover arrives or the preview toggle changes. The cover fills the folder body behind a thin folder-coloured rail that runs into the tab, making the picture large enough to tell two folders apart. With the toggle off, or with no cover, the body is a folder-coloured empty surface at the same size.

**Blocked by:** None (ticket 02 is done)

**Status:** done

- [x] Every directory row draws the folder shape whether "Folder preview" is on or off; the toggle only adds the cover picture
- [x] The shape is a tabbed silhouette: body 72x48, tab 32x16 proud of the body by 12px, 6px corner radius on the body, square join at the tab, subtle drop shadow
- [x] The cover fills the body inside a 2px folder-coloured rail, so the rail runs continuously from the body edge into the tab; the picture covers, never squashes
- [x] A directory with no cover, or the toggle off, draws the same silhouette with an empty folder-coloured body — no inner rectangle, no grey box, no image placeholder
- [x] The shape's flow box includes the tab: no part of the glyph overlaps the row above, the row's own text, or the row boundary, in the wide or the compact layout
- [x] A covered directory row, an uncovered directory row and an image-file row measure the same rendered height; switching the toggle changes no row's height
- [x] A covered directory row measures the same overall size with the toggle on and off
- [x] The compact 100px Preview column does not clip or overflow the shape
- [x] A cover that fails to load falls back to the identical empty folder shape at the same size
- [x] File rows keep their existing thumbnail markup and size
- [x] Component tests cover the identical footprint across the three row kinds, the unchanged height on toggle, the large cover, the empty body, and the failed-cover fallback
- [x] Full frontend test suite passes and `npm run lint:breakpoints` passes

## Comments

Geometry comes from the approved mock `design/folder-glyph-v2.html` (variant B, thin frame): body 72x48, tab 32x16 with 12px proud, silhouette 72x60, radius `0 6 6 6`, rail 2px, cover box 68x44. Measured on the mock: every row 66.6px, wide and compact, and the compact cell fits 72 + 8px padding.

Implementation note (deviation from the earlier "sits on the text baseline" line in the spec): the app loads Bootstrap, whose reboot sets `img { vertical-align: middle }`. A baseline-aligned inline-block therefore adds the line strut's descent below the 60px box and is 3.5px taller than the image thumbnail (64px vs 60.5px measured). The shape uses `vertical-align: middle` like the thumbnail, which makes its line-box contribution exactly the 60px silhouette and gives every row the same height. Measured in Chrome at 1200px and 412px: covered dir, uncovered dir and image file all 60.5px, and toggling changes none of them. The 60.5px (not the mock's 66.6px) is the app's real thumbnail row: a 60px border-box image plus the Material row's 0.5px border.
