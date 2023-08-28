package init

import (
	"gitea.suyono.dev/suyono/wingmate"
	"sync"
)

type Path interface {
	Path() string
}

type CronExactSpec interface {
	CronTimeSpec
	Value() uint8
}

type CronMultipleOccurrenceSpec interface {
	CronTimeSpec
	MultipleValues() []uint8
}

type CronTimeSpec interface {
	Type() wingmate.CronTimeType
}

type Cron interface {
	Minute() CronTimeSpec
	Hour() CronTimeSpec
	DayOfMonth() CronTimeSpec
	Month() CronTimeSpec
	DayOfWeek() CronTimeSpec
	Command() Path
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
	go i.sighandler(wg, signalTrigger, sighandlerExit)

	for _, s := range i.config.Services() {
		wg.Add(1)
		go func(p Path) {
			for {
				if err := i.service(wg, p); err != nil {
					wingmate.Log().Error().Msgf("starting service %s error %#v", p.Path(), err)
				}
			}
		}(s)
	}

	for _, c := range i.config.Cron() {
		wg.Add(1)
		go i.cron(wg, c)
	}
	wg.Wait()
}
