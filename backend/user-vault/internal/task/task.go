package task

import (
	"github.com/curtisnewbie/miso/middleware/task"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/vault"
)

func ScheduleTasks(rail miso.Rail) error {
	return task.ScheduleDistributedTask(
		miso.Job{
			Cron:                   miso.CronExprEveryXMin(30),
			Name:                   "LoadRoleAccessCacheTask",
			TriggeredOnBoostrapped: true,
			Run:                    vault.BatchLoadRoleAccessCache,
		},
		miso.Job{
			Cron:                   miso.CronExprEveryXMin(30),
			Name:                   "LoadPublicAccessCacheTask",
			TriggeredOnBoostrapped: true,
			Run:                    vault.LoadPublicAccessCache,
		},
	)
}
