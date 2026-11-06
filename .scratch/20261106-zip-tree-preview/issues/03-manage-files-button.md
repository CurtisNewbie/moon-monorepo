# 03: Open the Zip tree from the expanded row

**What to build:** A Zip file's expanded row on the manage-files page gains a Show tree button. The button appears only for `.zip` files. Clicking it requests the preview for that file and opens a dialog on success. On failure the dialog does not open and the error message appears in the existing snackbar.

The dialog renders the archive tree with the app's existing tree control. Folders and files show their names and icons, and files show the size formatted by the existing size formatter. Top-level nodes start expanded and deeper folders start collapsed. When the reply is truncated, the dialog says the tree was cut and names the cause. An archive with no entries shows an empty state. Names render as text.

All user-facing copy goes through the existing i18n mechanism, and the dialog is sized through the existing responsive dialog helper.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] The Show tree button appears in the expanded row only for files whose name ends in `.zip`.
- [x] Clicking the button requests the preview for that file key and opens the dialog with the returned tree.
- [x] The dialog renders folder and file nodes with names and icons, file sizes through the existing size formatter, and no size on folders.
- [x] Top-level nodes start expanded and deeper folders start collapsed.
- [x] A truncated reply shows the cut notice with its causes. A complete reply shows no notice.
- [x] An archive with no entries shows an empty state.
- [x] A failed request shows the error in the snackbar and opens no dialog.
- [x] Entry names render as text, never as markup.
- [x] Copy is added through the i18n dictionary and the dialog uses the responsive size helper.
- [x] Component tests cover rendering, sizes, truncation, empty state, and error handling. The manage-files spec covers button visibility and opening the dialog with the row's file key.
