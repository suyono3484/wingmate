package init

import (
	"os"
	"os/signal"
	"sync"

	"gitea.suyono.dev/suyono/wingmate"
	"golang.org/x/sys/unix"
)

func (i *Init) sighandler(wg *sync.WaitGroup, trigger chan<- any, selfExit <-chan any, sigchld chan<- os.Signal) {
	defer wg.Done()

	defer func() {
		wingmate.Log().Warn().Msg("signal handler: exiting")
	}()

	isOpen := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGINT, unix.SIGTERM, unix.SIGCHLD)

signal:
	for {
		select {
		case s := <-c:
			switch s {
			case unix.SIGTERM, unix.SIGINT:
				if isOpen {
					wingmate.Log().Info().Msg("initiating shutdown...")
					close(trigger)
					wg.Add(1)
					go i.signalPump(wg, selfExit)
					isOpen = false
				}
			case unix.SIGCHLD:
				select {
				case sigchld <- s:
				default:
				}
			}

		case <-selfExit:
			wingmate.Log().Warn().Msg("signal handler received completion flag")
			break signal
		}
	}
}
