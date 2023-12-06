package init

import (
	"errors"
	"sync"

	"golang.org/x/sys/unix"
)

func (i *Init) waiter(wg *sync.WaitGroup, runningFlag <-chan any, sigHandlerFlag chan<- any) {
	var (
		ws unix.WaitStatus
		// pid     int
		err     error
		running bool
	)

	running = true
wait:
	for {
		select {
		case <-runningFlag:
			running = false
		default:
		}

		if _, err = unix.Wait4(-1, &ws, 0, nil); err != nil {
			if errors.Is(err, unix.ECHILD) {
				if !running {
					close(sigHandlerFlag)
					break wait
				}
			}
		}
	}
}
