package main

import (
	"sync"

	"gitea.suyono.dev/suyono/wingmate/config"
	wminit "gitea.suyono.dev/suyono/wingmate/init"
	"gitea.suyono.dev/suyono/wingmate/task"
)

type wConfig struct {
	tasks    *task.Tasks
	config   *config.Config
	viperMtx *sync.Mutex
}

func (c *wConfig) Tasks() wminit.Tasks {
	return c.tasks
}

func (c *wConfig) Reload() error {
	//NOTE: for future use when reloading is possible
	return nil
}

func convert(cfg *config.Config) *wConfig {
	retval := &wConfig{
		tasks:    task.NewTasks(),
		config:   cfg,
		viperMtx: &sync.Mutex{},
	}

	for _, s := range cfg.Service {
		st := task.NewServiceTask(s.Name).SetCommand(s.Command...).SetEnv(s.Environ...)
		st.SetFlagSetsid(s.Setsid).SetWorkingDir(s.WorkingDir)
		st.SetUser(s.User).SetGroup(s.Group).SetStartSecs(s.StartSecs).SetPidFile(s.PidFile)
		st.SetConfig(cfg)
		retval.tasks.AddService(st)
	}

	for _, s := range cfg.ServiceV0 {
		retval.tasks.AddV0Service(s)
	}

	var schedule task.CronSchedule
	for _, c := range cfg.CronV0 {
		schedule = configToTaskCronSchedule(c.CronSchedule)
		retval.tasks.AddV0Cron(schedule, c.Command)
	}

	for _, c := range cfg.Cron {
		schedule = configToTaskCronSchedule(c.CronSchedule)

		ct := task.NewCronTask(c.Name).SetCommand(c.Command...).SetEnv(c.Environ...)
		ct.SetFlagSetsid(c.Setsid).SetWorkingDir(c.WorkingDir).SetUser(c.User).SetGroup(c.Group)
		ct.SetSchedule(schedule)
		ct.SetConfig(cfg)

		retval.tasks.AddCron(ct)
	}

	return retval
}

func configToTaskCronSchedule(cfgSchedule config.CronSchedule) (taskSchedule task.CronSchedule) {
	taskSchedule.Minute = configToTaskCronTimeSpec(cfgSchedule.Minute)
	taskSchedule.Hour = configToTaskCronTimeSpec(cfgSchedule.Hour)
	taskSchedule.DoM = configToTaskCronTimeSpec(cfgSchedule.DoM)
	taskSchedule.Month = configToTaskCronTimeSpec(cfgSchedule.Month)
	taskSchedule.DoW = configToTaskCronTimeSpec(cfgSchedule.DoW)

	return
}

func configToTaskCronTimeSpec(cfg config.CronTimeSpec) task.CronTimeSpec {
	switch v := cfg.(type) {
	case *config.SpecAny:
		return task.NewCronAnySpec()
	case *config.SpecExact:
		return task.NewCronExactSpec(v.Value())
	case *config.SpecMultiOccurrence:
		return task.NewCronMultiOccurrenceSpec(v.Values()...)
	}

	panic("invalid conversion")
}
