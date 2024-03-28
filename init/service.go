package init

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
)

const (
	serviceTag = "service"
)

func (i *Init) service(wg *sync.WaitGroup, task ServiceTask, exitFlag <-chan any) {
	defer wg.Done()

	var (
		err        error
		iwg        *sync.WaitGroup
		stderr     io.ReadCloser
		stdout     io.ReadCloser
		failStatus bool
		cmd        *exec.Cmd
	)

	defer func() {
		wingmate.Log().Info().Str(serviceTag, task.Name()).Msg("stopped")
	}()

service:
	for {
		failStatus = false
		if err = task.UtilDepCheck(); err != nil {
			wingmate.Log().Error().Str(serviceTag, task.Name()).Msgf("%+v", err)
			failStatus = true
			goto fail
		}
		cmd = exec.Command(task.Command(), task.Arguments()...)
		iwg = &sync.WaitGroup{}

		if stdout, err = cmd.StdoutPipe(); err != nil {
			wingmate.Log().Error().Str(serviceTag, task.Name()).Msgf("stdout pipe: %#v", err)
			failStatus = true
			goto fail
		}
		iwg.Add(1)
		go i.pipeReader(iwg, stdout, serviceTag, task.Name())

		if stderr, err = cmd.StderrPipe(); err != nil {
			wingmate.Log().Error().Str(serviceTag, task.Name()).Msgf("stderr pipe: %#v", err)
			_ = stdout.Close()
			failStatus = true
			goto fail
		}
		iwg.Add(1)
		go i.pipeReader(iwg, stderr, serviceTag, task.Name())

		if err = cmd.Start(); err != nil {
			wingmate.Log().Error().Msgf("starting service %s error %#v", task.Name(), err)
			failStatus = true
			_ = stdout.Close()
			_ = stderr.Close()
			iwg.Wait()
			goto fail
		}

		iwg.Wait()

		_ = cmd.Wait()

	fail:
		if failStatus {
			time.Sleep(time.Second)
			failStatus = false
		}
		select {
		case <-exitFlag:
			break service
		default:
		}
	}

}

func (i *Init) pipeReader(wg *sync.WaitGroup, pipe io.ReadCloser, tag, serviceName string) {
	defer wg.Done()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		wingmate.Log().Info().Str(tag, serviceName).Msg(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		wingmate.Log().Error().Str(tag, serviceName).Msgf("got error when reading pipe: %#v", err)
	}

	wingmate.Log().Info().Str(tag, serviceName).Msg("closing pipe")
}
