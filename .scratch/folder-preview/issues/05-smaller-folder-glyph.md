# 05: Smaller folder glyph

**What to build:** The folder shape drawn for every directory row shrinks from the 72x60 silhouette to a 54x44 one, so it stops dominating the Preview column and reads as subordinate to the rest of the row, while the cover picture inside it stays large enough to recognise a folder by. Row heights do not change: the cell keeps a fixed 60px slot and the glyph is centred inside it.

**Blocked by:** None (ticket 04 is done)

**Status:** ready-for-agent

- [ ] The folder body is 54x34 with a 5px corner radius on the body and a square join at the tab
- [ ] The tab is 24x10, proud of the body by 10px, with a 1px tuck at its base so no seam shows where it meets the body
- [ ] The cover fills the body inside the 2px folder-coloured rail: 50x30, `object-fit: cover`
- [ ] The Preview cell holds a fixed 60px slot and centres the glyph inside it, so a directory row's height does not depend on the glyph size
- [ ] A covered directory row, an uncovered directory row and an image-file row still measure the same row height, wide and compact
- [ ] No part of the glyph, tab included, overlaps the row above, the row's own text, or the row boundary
- [ ] The empty-body state, the failed-cover fallback and the compact 100px column behave exactly as before
- [ ] File rows keep their existing markup and size
- [ ] Component tests cover the new measured footprint (54x44 body+tab, 50x30 cover) and the unchanged row heights
- [ ] Full frontend test suite passes and `npm run lint:breakpoints` passes

## Comments

Size chosen from `design/folder-glyph-v3.html` (measured there: all rows 61px for S/M/L, glyph 54x34 body, 54x44 silhouette, cover 50x30; 1.14x the width of the "Preview" header, versus 1.53x for the shipped 72-wide shape).
