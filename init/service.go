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

func (i *Init) service(wg *sync.WaitGroup, path Path, exitFlag <-chan any) {
	defer wg.Done()

	var (
		err        error
		iwg        *sync.WaitGroup
		stderr     io.ReadCloser
		stdout     io.ReadCloser
		failStatus bool
	)

	defer func() {
		wingmate.Log().Info().Str(serviceTag, path.Path()).Msg("stopped")
	}()

service:
	for {
		failStatus = false
		cmd := exec.Command(path.Path())
		iwg = &sync.WaitGroup{}

		if stdout, err = cmd.StdoutPipe(); err != nil {
			wingmate.Log().Error().Str(serviceTag, path.Path()).Msgf("stdout pipe: %#v", err)
			failStatus = true
			goto fail
		}
		iwg.Add(1)
		go i.pipeReader(iwg, stdout, serviceTag, path.Path())

		if stderr, err = cmd.StderrPipe(); err != nil {
			wingmate.Log().Error().Str(serviceTag, path.Path()).Msgf("stderr pipe: %#v", err)
			_ = stdout.Close()
			failStatus = true
			goto fail
		}
		iwg.Add(1)
		go i.pipeReader(iwg, stderr, serviceTag, path.Path())

		if err = cmd.Start(); err != nil {
			wingmate.Log().Error().Msgf("starting service %s error %#v", path.Path(), err)
			failStatus = true
			_ = stdout.Close()
			_ = stderr.Close()
			iwg.Wait()
			goto fail
		}

		iwg.Wait()

		if err = cmd.Wait(); err != nil {
			wingmate.Log().Error().Str(serviceTag, path.Path()).Msgf("got error when waiting: %+v", err)
		}
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
