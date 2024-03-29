package init

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
)

const (
	cronTag = "cron"
)

func (i *Init) cron(wg *sync.WaitGroup, cron CronTask, exitFlag <-chan any) {
	defer wg.Done()

	var (
		iwg    *sync.WaitGroup
		err    error
		stdout io.ReadCloser
		stderr io.ReadCloser
		cmd    *exec.Cmd
	)

	ticker := time.NewTicker(time.Second * 30)
cron:
	for {
		if cron.TimeToRun(time.Now()) {
			wingmate.Log().Info().Str(cronTag, cron.Name()).Msg("executing")
			if err = cron.UtilDepCheck(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Name()).Msgf("%+v", err)
				goto fail
			}
			cmd = exec.Command(cron.Command(), cron.Arguments()...)
			cmd.Env = os.Environ()
			if cron.EnvLen() > 0 {
				cmd.Env = append(cmd.Env, cron.Environ()...)
			}

			if len(cron.WorkingDir()) > 0 {
				cmd.Dir = cron.WorkingDir()
			}

			iwg = &sync.WaitGroup{}

			if stdout, err = cmd.StdoutPipe(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Name()).Msgf("stdout pipe: %+v", err)
				goto fail
			}

			if stderr, err = cmd.StderrPipe(); err != nil {
				wingmate.Log().Error().Str(cronTag, cron.Name()).Msgf("stderr pipe: %+v", err)
				_ = stdout.Close()
				goto fail
			}

			iwg.Add(1)
			go i.pipeReader(iwg, stdout, cronTag, cron.Name())

			iwg.Add(1)
			go i.pipeReader(iwg, stderr, cronTag, cron.Name())

			if err := cmd.Start(); err != nil {
				wingmate.Log().Error().Msgf("starting cron %s error %+v", cron.Name(), err)
				_ = stdout.Close()
				_ = stderr.Close()
				iwg.Wait()
				goto fail
			}

			iwg.Wait()

			_ = cmd.Wait()
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
