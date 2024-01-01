package task

import (
	wminit "gitea.suyono.dev/suyono/wingmate/init"
)

type Tasks struct {
	services []wminit.ServiceTask
	crones   []wminit.CronTask
}

func NewTasks() *Tasks {
	return &Tasks{
		services: make([]wminit.ServiceTask, 0),
		crones:   make([]wminit.CronTask, 0),
	}
}

func (ts *Tasks) AddV0Service(path string) {
	ts.services = append(ts.services, &Task{
		name:    path,
		command: []string{path},
	})
}

func (ts *Tasks) AddV0Cron(schedule CronSchedule, path string) {
	ts.crones = append(ts.crones, &Cron{
		CronSchedule: schedule,
		name:         path,
		command:      []string{path},
		hasRun:       false,
	})
}

func (ts *Tasks) List() []wminit.Task {
	retval := make([]wminit.Task, 0, len(ts.services)+len(ts.crones))
	for _, s := range ts.services {
		retval = append(retval, s.(wminit.Task))
	}
	for _, c := range ts.crones {
		retval = append(retval, c.(wminit.Task))
	}
	return retval
}

func (ts *Tasks) Services() []wminit.ServiceTask {
	return ts.services
}

func (ts *Tasks) Crones() []wminit.CronTask {
	return ts.crones
}

func (ts *Tasks) Get(name string) (wminit.Task, error) {
	panic("not implemented")
	return nil, nil
}

type Task struct {
	name    string
	command []string
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) Command() []string {
	retval := make([]string, len(t.command))
	copy(retval, t.command)
	return retval
}

func (t *Task) Environ() []string {
	panic("not implemented")
	return nil
}

func (t *Task) Setsid() bool {
	panic("not implemented")
	return false
}

func (t *Task) UserGroup() wminit.UserGroup {
	panic("not implemented")
	return nil
}

func (t *Task) Background() bool {
	panic("not implemented")
	return false
}

func (t *Task) WorkingDir() string {
	panic("not implemented")
	return ""
}

func (t *Task) Status() wminit.TaskStatus {
	panic("not implemented")
	return nil
}

func (t *Task) AutoStart() bool {
	panic("not implemented")
	return false
}

func (t *Task) AutoRestart() bool {
	panic("not implemented")
	return false
}
