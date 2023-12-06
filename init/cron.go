package init

import (
	"os/exec"
	"sync"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
)

func (i *Init) cron(wg *sync.WaitGroup, cron Cron, exitFlag <-chan any) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second * 20)
cron:
	for {
		if cron.TimeToRun(time.Now()) {
			cmd := exec.Command(cron.Command().Path())
			if err := cmd.Run(); err != nil {
				wingmate.Log().Error().Msgf("running cron %s error %#v", cron.Command().Path(), err)
			}
		}

		select {
		case <-exitFlag:
			ticker.Stop()
			break cron
		case <-ticker.C:
		}
	}
}
