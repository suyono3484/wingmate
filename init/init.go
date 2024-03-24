package init

import (
	"os"
	"sync"
	"time"
)

type Tasks interface {
	List() []Task
	Services() []ServiceTask
	Crones() []CronTask
	Get(string) (Task, error)
}

type UserGroup interface {
	String() string
	IsSet() bool
}

type TaskStatus interface {
}

type Task interface {
	Name() string
	Command() string
	Arguments() []string
	Environ() []string
	Setsid() bool
	UserGroup() UserGroup
	WorkingDir() string
	Status() TaskStatus
}

type CronTask interface {
	Task
	TimeToRun(time.Time) bool
}

type ServiceTask interface {
	Task
	Background() bool //NOTE: implies using wmpidproxy
	PidFile() string  //NOTE: implies using wmpidproxy
	StartSecs() uint
	AutoStart() bool
	AutoRestart() bool
}

type Config interface {
	Tasks() Tasks
}

type Init struct {
	config Config
}

func NewInit(config Config) *Init {
	return &Init{
		config: config,
	}
}

func (i *Init) Start() {
	var (
		wg             *sync.WaitGroup
		signalTrigger  chan any
		sighandlerExit chan any
		sigchld        chan os.Signal
	)

	signalTrigger = make(chan any)
	sighandlerExit = make(chan any)
	sigchld = make(chan os.Signal, 1)

	wg = &sync.WaitGroup{}
	wg.Add(1)
	go i.waiter(wg, signalTrigger, sighandlerExit, sigchld)

	wg.Add(1)
	go i.sighandler(wg, signalTrigger, sighandlerExit, sigchld)

	for _, s := range i.config.Tasks().Services() {
		wg.Add(1)
		go i.service(wg, s, signalTrigger)
	}

	for _, c := range i.config.Tasks().Crones() {
		wg.Add(1)
		go i.cron(wg, c, signalTrigger)
	}
	wg.Wait()
}
