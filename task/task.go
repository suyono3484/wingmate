package task

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"gitea.suyono.dev/suyono/wingmate"

	wminit "gitea.suyono.dev/suyono/wingmate/init"
)

type config interface {
	WMPidProxyPath() string
	WMPidProxyCheckVersion() error
	WMExecPath() string
	WMExecCheckVersion() error
}

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
	ts.AddCron(NewCronTask(path)).SetCommand(path).SetSchedule("", schedule)
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
	//TODO: implement me!
	panic("not implemented")
	return nil, nil
}

type ServiceTask struct {
	name       string
	command    []string
	cmdLine    []string
	environ    []string
	setsid     bool
	background bool
	workingDir string
	startSecs  uint
	pidFile    string
	config     config
	userGroup
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

func (t *ServiceTask) SetStartSecs(secs uint) *ServiceTask {
	t.startSecs = secs
	return t
}

func (t *ServiceTask) SetPidFile(path string) *ServiceTask {
	t.pidFile = path
	if len(path) > 0 {
		t.background = true
	} else {
		t.background = false
	}
	return t
}

func (t *ServiceTask) SetConfig(config config) *ServiceTask {
	t.config = config
	return t
}

func (t *ServiceTask) Equals(another *ServiceTask) bool {
	if another == nil {
		return false
	}

	type toCompare struct {
		Name        string
		Command     string
		Arguments   []string
		Environ     []string
		Setsid      bool
		UserGroup   string
		WorkingDir  string
		PidFile     string
		StartSecs   uint
		AutoStart   bool
		AutoRestart bool
	}

	cmpStruct := func(p *ServiceTask) ([]byte, error) {
		s := &toCompare{
			Name:        p.Name(),
			Command:     p.Command(),
			Arguments:   p.Arguments(),
			Environ:     p.Environ(),
			Setsid:      p.Setsid(),
			UserGroup:   p.UserGroup().String(),
			WorkingDir:  p.WorkingDir(),
			PidFile:     p.PidFile(),
			StartSecs:   p.StartSecs(),
			AutoStart:   p.AutoStart(),
			AutoRestart: p.AutoRestart(),
		}

		return json.Marshal(s)
	}

	var (
		err                error
		ours, theirs       []byte
		ourHash, theirHash [sha256.Size]byte
	)

	if ours, err = cmpStruct(t); err != nil {
		wingmate.Log().Error().Msgf("task equals: %+v", err)
		return false
	}
	ourHash = sha256.Sum256(ours)

	if theirs, err = cmpStruct(another); err != nil {
		wingmate.Log().Error().Msgf("task equals: %+v", err)
		return false
	}
	theirHash = sha256.Sum256(theirs)

	for i := 0; i < sha256.Size; i++ {
		if ourHash[i] != theirHash[i] {
			return false
		}
	}

	return true
}

func (t *ServiceTask) Validate() error {
	// call this function for validate the field
	return validate( /* input the validators here */ )
}

func (t *ServiceTask) Name() string {
	return t.name
}

func (t *ServiceTask) prepareCommandLine() []string {
	if len(t.cmdLine) > 0 {
		return t.cmdLine
	}

	t.cmdLine = make([]string, 0)
	if t.background {
		t.cmdLine = append(t.cmdLine, t.config.WMPidProxyPath(), "--pid-file", t.pidFile, "--")
	}

	if t.setsid || t.UserGroup().IsSet() {
		t.cmdLine = append(t.cmdLine, t.config.WMExecPath())

		if t.setsid {
			t.cmdLine = append(t.cmdLine, "--setsid")
		}

		if t.UserGroup().IsSet() {
			t.cmdLine = append(t.cmdLine, "--user", t.UserGroup().String())
		}

		t.cmdLine = append(t.cmdLine, "--")
	}

	t.cmdLine = append(t.cmdLine, t.command...)

	return t.cmdLine
}

func (t *ServiceTask) UtilDepCheck() error {
	t.cmdLine = make([]string, 0)
	if t.background {
		if err := t.config.WMPidProxyCheckVersion(); err != nil {
			return fmt.Errorf("utility dependency check: %w", err)
		}

		t.cmdLine = append(t.cmdLine, t.config.WMPidProxyPath(), "--pid-file", t.pidFile, "--")
	}

	if t.setsid || t.UserGroup().IsSet() {
		if err := t.config.WMExecCheckVersion(); err != nil {
			return fmt.Errorf("utility dependency check: %w", err)
		}

		t.cmdLine = append(t.cmdLine, t.config.WMExecPath())

		if t.setsid {
			t.cmdLine = append(t.cmdLine, "--setsid")
		}

		if t.UserGroup().IsSet() {
			t.cmdLine = append(t.cmdLine, "--user", t.UserGroup().String())
		}

		t.cmdLine = append(t.cmdLine, "--")
	}

	t.cmdLine = append(t.cmdLine, t.command...)

	return nil
}

func (t *ServiceTask) Command() string {
	return t.cmdLine[0]
}

func (t *ServiceTask) Arguments() []string {
	if len(t.cmdLine) == 1 {
		return nil
	}

	retval := make([]string, len(t.cmdLine)-1)
	copy(retval, t.cmdLine[1:])

	return retval
}

func (t *ServiceTask) EnvLen() int {
	return len(t.environ)
}

func (t *ServiceTask) Environ() []string {
	retval := make([]string, len(t.environ))
	copy(retval, t.environ)
	return retval
}

func (t *ServiceTask) Setsid() bool {
	return t.setsid
}

func (t *ServiceTask) UserGroup() wminit.UserGroup {
	return &(t.userGroup)
}

func (t *ServiceTask) Background() bool {
	return t.background
}

func (t *ServiceTask) WorkingDir() string {
	return t.workingDir
}

func (t *ServiceTask) Status() wminit.TaskStatus {
	//TODO: implement me!
	panic("not implemented")
	return nil
}

func (t *ServiceTask) AutoStart() bool {
	//TODO: implement me!
	panic("not implemented")
	return false
}

func (t *ServiceTask) AutoRestart() bool {
	//TODO: implement me!
	panic("not implemented")
	return false
}

func (t *ServiceTask) StartSecs() uint {
	return t.startSecs
}

func (t *ServiceTask) PidFile() string {
	return t.pidFile
}

type userGroup struct {
	user  string
	group string
}

func (ug *userGroup) IsSet() bool {
	return len(ug.user) > 0 || len(ug.group) > 0
}

func (ug *userGroup) String() string {
	if len(ug.group) > 0 {
		return fmt.Sprintf("%s:%s", ug.user, ug.group)
	}

	return ug.user
}

func validate(validators ...func() error) error {
	var err error
	for _, v := range validators {
		if err = v(); err != nil {
			return err
		}
	}
	return nil
}
