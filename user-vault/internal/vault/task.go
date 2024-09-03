package vault

import (
	"github.com/curtisnewbie/miso/miso"
)

func ScheduleTasks(rail miso.Rail) error {
	// distributed tasks
	var err error = miso.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadRoleResCacheTask",
		TriggeredOnBoostrapped: true,
		Run:                    LoadRoleResCache,
	})
	if err != nil {
		return err
	}
	err = miso.ScheduleDistributedTask(miso.Job{
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
