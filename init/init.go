package init

import (
	"sync"
	"time"
)

type Path interface {
	Path() string
}

type Cron interface {
	Command() Path
	TimeToRun(time.Time) bool
}

type Config interface {
	Services() []Path
	Cron() []Cron
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
	)

	signalTrigger = make(chan any)
	sighandlerExit = make(chan any)

	wg = &sync.WaitGroup{}
	wg.Add(1)
	go i.waiter(wg, signalTrigger, sighandlerExit)

	wg.Add(1)
	go i.sighandler(wg, signalTrigger, sighandlerExit)

	for _, s := range i.config.Services() {
		wg.Add(1)
		go i.service(wg, s, signalTrigger)
	}

	for _, c := range i.config.Cron() {
		wg.Add(1)
		go i.cron(wg, c, signalTrigger)
	}
	wg.Wait()
}
