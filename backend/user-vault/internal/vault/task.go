package vault

import (
	"github.com/curtisnewbie/miso/middleware/task"
	"github.com/curtisnewbie/miso/miso"
)

func ScheduleTasks(rail miso.Rail) error {
	// distributed tasks
	var err error = task.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadRoleResCacheTask",
		TriggeredOnBoostrapped: true,
		Run:                    LoadRoleResCache,
	})
	if err != nil {
		return err
	}
	err = task.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadPathResCacheTask",
		TriggeredOnBoostrapped: true,
		Run:                    LoadPathResCache,
	})
	if err != nil {
		return err
	}
	return nil
}
