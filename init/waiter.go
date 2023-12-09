package init

import (
	"errors"
	"sync"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
	"golang.org/x/sys/unix"
)

func (i *Init) waiter(wg *sync.WaitGroup, runningFlag <-chan any, sigHandlerFlag chan<- any) {
	var (
		ws unix.WaitStatus
		// pid     int
		err     error
		running bool
		flagged bool
	)
	defer wg.Done()

	defer func() {
		wingmate.Log().Info().Msg("waiter exiting...")
	}()

	running = true
	flagged = true
wait:
	for {
		select {
		case <-runningFlag:
			wingmate.Log().Info().Msg("waiter received shutdown signal...")
			running = false
		default:
		}

		if _, err = unix.Wait4(-1, &ws, 0, nil); err != nil {
			if errors.Is(err, unix.ECHILD) {
				if !running {
					if flagged {
						close(sigHandlerFlag)
						flagged = false
						wingmate.Log().Warn().Msg("waiter: inner flag")
					}
					wingmate.Log().Warn().Msg("waiter: no child left")
					break wait
				}
			}

			wingmate.Log().Warn().Msgf("Wait4 returns error: %+v", err)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
