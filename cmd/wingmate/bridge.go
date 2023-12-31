package main

import (
	"gitea.suyono.dev/suyono/wingmate/config"
	wminit "gitea.suyono.dev/suyono/wingmate/init"
	"gitea.suyono.dev/suyono/wingmate/task"
)

type wConfig struct {
	tasks *task.Tasks
}

func (c *wConfig) Tasks() wminit.Tasks {
	return c.tasks
}

func convert(cfg *config.Config) *wConfig {
	retval := &wConfig{
		tasks: task.NewTasks(),
	}

	for _, s := range cfg.ServicePaths {
		retval.tasks.AddV0Service(s)

	}

	var schedule task.CronSchedule
	for _, c := range cfg.Cron {
		schedule.Minute = convertSchedule(c.Minute)
		schedule.Hour = convertSchedule(c.Hour)
		schedule.DoM = convertSchedule(c.DoM)
		schedule.Month = convertSchedule(c.Month)
		schedule.DoW = convertSchedule(c.DoW)

		retval.tasks.AddV0Cron(schedule, c.Command)
	}

	return retval
}

func convertSchedule(cfg config.CronTimeSpec) task.CronTimeSpec {
	switch v := cfg.(type) {
	case *config.SpecAny:
		return task.NewCronAnySpec()
	case *config.SpecExact:
		return task.NewCronExactSpec(v.Value())
	case *config.SpecMultiOccurrence:
		return task.NewCronMultiOccurrenceSpec(v.Values()...)
	}

	panic("invalid conversion")
	return nil
}
