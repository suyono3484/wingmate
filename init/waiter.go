package init

import (
	"errors"
	"os"
	"sync"

	"gitea.suyono.dev/suyono/wingmate"
	"golang.org/x/sys/unix"
)

func (i *Init) waiter(wg *sync.WaitGroup, runningFlag <-chan any, sigHandlerFlag chan<- any, sigchld <-chan os.Signal) {
	var (
		ws               unix.WaitStatus
		err              error
		running          bool
		flagged          bool
		waitingForSignal bool
	)
	defer wg.Done()

	defer func() {
		wingmate.Log().Info().Msg("waiter exiting...")
	}()

	running = true
	flagged = true
	waitingForSignal = true
wait:
	for {
		if running {
			if waitingForSignal {
				select {
				case <-runningFlag:
					wingmate.Log().Info().Msg("waiter received shutdown signal...")
					running = false
				case <-sigchld:
					waitingForSignal = false
				}

			} else {
				select {
				case <-runningFlag:
					wingmate.Log().Info().Msg("waiter received shutdown signal...")
					running = false
				default:
				}

			}
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
			waitingForSignal = true
		}
	}
}
