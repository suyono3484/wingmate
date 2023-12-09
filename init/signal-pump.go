package init

import (
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

type status int

const (
	triggered status = iota
	expired
)

func (i *Init) signalPump(wg *sync.WaitGroup, selfExit <-chan any) {
	if seStatus := i.sigTermPump(time.Now(), selfExit); seStatus == triggered {
		return
	}

	i.sigKillPump(time.Now(), selfExit)
}

func (i *Init) sigKillPump(startTime time.Time, selfExit <-chan any) {
	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()

	for time.Since(startTime) < time.Second {
		_ = unix.Kill(-1, unix.SIGKILL)

		select {
		case <-t.C:
		case <-selfExit:
			return
		}
	}
}

func (i *Init) sigTermPump(startTime time.Time, selfExit <-chan any) status {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for time.Since(startTime) < time.Duration(time.Second*4) {
		_ = unix.Kill(-1, unix.SIGTERM)

		select {
		case <-t.C:
		case <-selfExit:
			return triggered
		}
	}

	return expired
}
