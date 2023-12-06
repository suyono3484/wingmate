package init

import (
	"os"
	"os/signal"
	"sync"

	"golang.org/x/sys/unix"
)

func (i *Init) sighandler(wg *sync.WaitGroup, trigger chan<- any, selfExit <-chan any) {
	defer wg.Wait()

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
					close(trigger)
					isOpen = false
				}
			case unix.SIGCHLD:
				// do nothing
			}

		case <-selfExit:
			break signal
		}
	}
}
