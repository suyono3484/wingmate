package init

import "sync"

func (i *Init) cron(wg *sync.WaitGroup, cron Cron) {
	defer wg.Done()
}
