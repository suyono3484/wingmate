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
	ts.AddService(NewServiceTask(path)).SetCommand(path)
}

func (ts *Tasks) AddService(serviceTask *ServiceTask) *ServiceTask {
	ts.services = append(ts.services, serviceTask)
	return serviceTask
}

func (ts *Tasks) AddV0Cron(schedule CronSchedule, path string) {
	ts.AddCron(NewCronTask(path)).SetCommand(path).SetSchedule(schedule)
}

func (ts *Tasks) AddCron(cronTask *CronTask) *CronTask {
	ts.crones = append(ts.crones, cronTask)
	return cronTask
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

type ServiceTask struct {
	name       string
	command    []string
	environ    []string
	setsid     bool
	user       string
	group      string
	background bool
	workingDir string
}

func NewServiceTask(name string) *ServiceTask {
	return &ServiceTask{
		name: name,
	}
}

func (t *ServiceTask) SetCommand(cmds ...string) *ServiceTask {
	t.command = make([]string, len(cmds))
	copy(t.command, cmds)
	return t
}

func (t *ServiceTask) SetEnv(envs ...string) *ServiceTask {
	t.environ = make([]string, len(envs))
	copy(t.environ, envs)
	return t
}

func (t *ServiceTask) SetFlagSetsid(flag bool) *ServiceTask {
	t.setsid = flag
	return t
}

func (t *ServiceTask) SetFlagBackground(flag bool) *ServiceTask {
	t.background = flag
	return t
}

func (t *ServiceTask) SetWorkingDir(path string) *ServiceTask {
	t.workingDir = path
	return t
}

func (t *ServiceTask) SetUser(user string) *ServiceTask {
	t.user = user
	return t
}

func (t *ServiceTask) SetGroup(group string) *ServiceTask {
	t.group = group
	return t
}

func (t *ServiceTask) Name() string {
	return t.name
}

func (t *ServiceTask) Command() []string {
	retval := make([]string, len(t.command))
	copy(retval, t.command)
	return retval
}

func (t *ServiceTask) Environ() []string {
	panic("not implemented")
	return nil
}

func (t *ServiceTask) Setsid() bool {
	panic("not implemented")
	return false
}

func (t *ServiceTask) UserGroup() wminit.UserGroup {
	panic("not implemented")
	return nil
}

func (t *ServiceTask) Background() bool {
	panic("not implemented")
	return false
}

func (t *ServiceTask) WorkingDir() string {
	panic("not implemented")
	return ""
}

func (t *ServiceTask) Status() wminit.TaskStatus {
	panic("not implemented")
	return nil
}

func (t *ServiceTask) AutoStart() bool {
	panic("not implemented")
	return false
}

func (t *ServiceTask) AutoRestart() bool {
	panic("not implemented")
	return false
}
