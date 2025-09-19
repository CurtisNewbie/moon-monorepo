package vault

import (
	"github.com/curtisnewbie/miso/middleware/task"
	"github.com/curtisnewbie/miso/miso"
)

func ScheduleTasks(rail miso.Rail) error {
	err := task.ScheduleDistributedTask(miso.Job{
		Cron: "*/30 * * * *",

		Name:                   "LoadRoleAccessCacheTask",
		TriggeredOnBoostrapped: true,
		Run:                    BatchLoadRoleAccessCache,
	})
	if err != nil {
		return err
	}
	err = task.ScheduleDistributedTask(miso.Job{
		Cron: "*/30 * * * *",

		Name:                   "LoadPublicAccessCacheTask",
		TriggeredOnBoostrapped: true,
		Run:                    LoadPublicAccessCache,
	})
	if err != nil {
		return err
	}
	return nil
}
