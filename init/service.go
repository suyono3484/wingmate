package init

import (
	"os/exec"
	"sync"

	"gitea.suyono.dev/suyono/wingmate"
)

func (i *Init) service(wg *sync.WaitGroup, path Path, exitFlag <-chan any) {
	defer wg.Done()

	var (
		err error
	)

service:
	for {
		cmd := exec.Command(path.Path())
		if err = cmd.Run(); err != nil {
			wingmate.Log().Error().Msgf("starting service %s error %#v", path.Path(), err)
		}

		select {
		case <-exitFlag:
			break service
		default:
		}
	}

}
