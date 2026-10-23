# 05: Size sort selectable without setting the type filter

**What to build:** Size sits in the sort dropdown at all times, so a user can pick it in one action from any type filter. Selecting size sets the type filter to File and loads the size-ordered list; the user no longer has to choose File first. Picking size while already on File behaves as before.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] The size option appears in the sort dropdown under All, File, and Directory
- [ ] Selecting size under All or Directory sets the type filter to File and loads the size-ordered list in the same action
- [ ] Selecting size under File behaves as before
- [ ] Changing the type filter away from File still resets the sort to the stored or default order
- [ ] Size is still never sent to the order-by preference endpoint and still records no Resume Cursor
- [ ] Component tests replace the old files-only visibility assertions with the always-offered behavior and the one-action coercion
