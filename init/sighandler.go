package init

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"sync"
)

func (i *Init) sighandler(wg *sync.WaitGroup, trigger chan<- any, selfExit <-chan any) {
	defer wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGINT, unix.SIGTERM, unix.SIGCHLD)

signal:
	for {
		select {
		case s := <-c:
			switch s {
			case unix.SIGTERM, unix.SIGINT:
				close(trigger)
			case unix.SIGCHLD:
				// do nothing
			}

		case <-selfExit:
			break signal
		}
	}
}
