package config

import "github.com/curtisnewbie/miso/miso"

// misoconfig-section: mini-fstore configuration
const (
	// misoconfig-prop: Storage Directory | ./storage
	PropStorageDir = "fstore.storage.dir"

	// misoconfig-prop: Trash Directory | ./trash
	PropTrashDir = "fstore.trash.dir"

	// misoconfig-prop: Temporary directory | /tmp
	PropTempDir = "fstore.tmp.dir"

	// misoconfig-prop: Strategy used to 'physically' delete files, there are two types of strategies available: direct / trash. When using 'direct' strategy, files are deleted directly. When using 'trash' strategy, files are moved into the trash directory. | `"trash"`
	PropPDelStrategy = "fstore.pdelete.strategy"

	// misoconfig-prop: Enable endpoints for mini-fstore file backup, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup). | false
	PropBackupEnabled = "fstore.backup.enabled"

	// misoconfig-prop: Secret for backup endpoints authorization, see [fstore_backup](https://github.com/curtisnewbie/fstore_backup). |
	PropBackupAuthSecret = "fstore.backup.secret"

	// misoconfig-prop: Enable dry-run mode for StanitizeStorageTask | false
	PropSanitizeStorageTaskDryRun = "task.sanitize-storage-task.dry-run"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropStorageDir, "./storage")
	miso.SetDefProp(PropTrashDir, "./trash")
	miso.SetDefProp(PropTempDir, "/tmp")
	miso.SetDefProp(PropPDelStrategy, "trash")
	miso.SetDefProp(PropBackupEnabled, false)
	miso.SetDefProp(PropSanitizeStorageTaskDryRun, false)
}

// misoconfig-default-end
