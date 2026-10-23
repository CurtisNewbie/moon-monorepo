# 03: Size option in the manage-files UI

**What to build:** The manage-files page offers size in the sort dropdown only while the File type filter is selected. Choosing size while the filter is All or Directory switches the filter to File and loads the size-ordered list. Switching the filter away from File resets the sort to the directory's stored or default order. Size is never uploaded as the directory's Order-By Preference and records no Resume Cursor. The label is translated in every locale the page already ships.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] The size option appears in the sort dropdown only when the type filter is File
- [ ] Choosing size under All or Directory sets the type filter to File and fetches the size-ordered list
- [ ] Switching the type filter away from File while size is active resets the sort to the stored or default order and refetches
- [ ] Size is never sent to the order-by preference endpoint
- [ ] Browsing a size-ordered listing records no Resume Cursor and triggers no resume prompt
- [ ] The size label is present in every locale the page uses
- [ ] Component tests cover the dropdown visibility, the type coercion, the reset on type change, and the skipped preference upload
