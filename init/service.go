package init

import (
	"os/exec"
	"sync"
)

func (i *Init) service(wg *sync.WaitGroup, path Path) error {
	defer wg.Done()

	var (
		err error
	)

	cmd := exec.Command(path.Path())
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
