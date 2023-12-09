package init

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
)

const (
	cronTag = "cron"
)

func (i *Init) cron(wg *sync.WaitGroup, cron Cron, exitFlag <-chan any) {
	defer wg.Done()

	var (
		iwg    *sync.WaitGroup
		err    error
		stdout io.ReadCloser
		stderr io.ReadCloser
	)

	ticker := time.NewTicker(time.Second * 30)
cron:
	for {
		if cron.TimeToRun(time.Now()) {
			wingmate.Log().Info().Str(cronTag, cron.Command().Path()).Msg("executing")
			cmd := exec.Command(cron.Command().Path())
			iwg = &sync.WaitGroup{}

			if stdout, err = cmd.StdoutPipe(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Command().Path()).Msgf("stdout pipe: %+v", err)
				goto fail
			}

			if stderr, err = cmd.StderrPipe(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Command().Path()).Msgf("stderr pipe: %+v", err)
				_ = stdout.Close()
				goto fail
			}

			iwg.Add(1)
			go i.pipeReader(iwg, stdout, cronTag, cron.Command().Path())

			iwg.Add(1)
			go i.pipeReader(iwg, stderr, cronTag, cron.Command().Path())

			if err := cmd.Start(); err != nil {
				wingmate.Log().Error().Msgf("starting cron %s error %+v", cron.Command().Path(), err)
				_ = stdout.Close()
				_ = stderr.Close()
				iwg.Wait()
				goto fail
			}

			iwg.Wait()

			if err = cmd.Wait(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Command().Path()).Msgf("got error when waiting: %+v", err)
			}
		}

	fail:
		select {
		case <-exitFlag:
			ticker.Stop()
			break cron
		case <-ticker.C:
		}
	}
}
