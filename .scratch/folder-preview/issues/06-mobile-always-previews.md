# 06: Mobile always previews folders

**What to build:** On a phone the compact list always draws folder covers, and the "Folder Preview" checkbox is not offered there. The wide layout keeps the checkbox and its remembered choice; the compact layout ignores that choice and behaves as if it were on.

**Blocked by:** None (tickets 04 and 05 are done)

**Status:** done

- [x] In the compact layout a directory row with a known cover draws the cover even when the stored preference is off
- [x] In the compact layout the "Folder Preview" checkbox is not rendered
- [x] In the wide layout the checkbox still renders beside the tooltip checkbox, still defaults off, and still remembers the choice
- [x] In the wide layout turning the checkbox off still clears the covers
- [x] The mobile-only `order: -1` rule for the checkbox is removed along with the control
- [x] A directory with no cover still draws the empty folder body on a phone
- [x] Row heights are unchanged in both layouts
- [x] Component tests cover: no checkbox in the compact layout, covers drawn on a phone with the preference off, covers withheld on a phone with no cover known, and the wide layout unchanged
- [x] Full frontend test suite passes and `npm run lint:breakpoints` passes

## Comments

The compact strip scrolls sideways and the checkbox sat last in it, so it was reachable only by scrolling. Rather than move it, the compact layout drops the control and always previews: the cover is the reason the switch exists, and a phone has no room for a dead one.

Implemented: `showDirCover` returns true for a directory with a known cover when `this.folderPreview || this.env.isMobile()`; the checkbox carries `*ngIf="!env.isMobile()"` like the tooltip one; the `.folder-preview-toggle` class and its mobile `order: -1` rule are gone. The compact checkbox test was replaced by "does not offer the checkbox in the compact layout", and the compact footprint test no longer ticks anything — a phone previews without a switch.
