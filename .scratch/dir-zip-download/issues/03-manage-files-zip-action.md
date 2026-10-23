# 03: manage-files: expanded panel zip action

**What to build:** Opening a directory row in the manage-files list shows a zip download action in the expanded panel's directory actions, beside rename, comic toggle and truncate. Clicking it starts the download of that directory's archive. The action appears for directory rows only, its tooltip comes from the translation catalog, and file rows keep exactly the download they have today.

**Blocked by:** 02

**Status:** done

- [x] A directory row's expanded panel offers a zip download action with a translated tooltip.
- [x] File rows show no zip action and keep their existing single-file download unchanged.
- [x] Clicking the action asks vfm for the directory's bundle key and opens the archive route with it.
- [x] The action carries no literal user-facing string; the tooltip resolves through the translation catalog.
- [x] A failed key request surfaces through the app's existing error handling rather than a silent no-op.
