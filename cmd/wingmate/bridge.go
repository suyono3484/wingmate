package main

import (
	"time"

	"gitea.suyono.dev/suyono/wingmate/config"
	wminit "gitea.suyono.dev/suyono/wingmate/init"
)

type wPath struct {
	path string
}

func (p wPath) Path() string {
	return p.path
}

type wConfig struct {
	services []wminit.Path
	cron     []wminit.Cron
}

func (c wConfig) Services() []wminit.Path {
	return c.services
}

func (c wConfig) Cron() []wminit.Cron {
	return c.cron
}

type wCron struct {
	iCron *config.Cron
}

func (c wCron) TimeToRun(now time.Time) bool {
	return c.iCron.TimeToRun(now)
}

func (c wCron) Command() wminit.Path {
	return wPath{
		path: c.iCron.Command(),
	}
}

func convert(cfg *config.Config) wConfig {
	retval := wConfig{
		services: make([]wminit.Path, 0, len(cfg.ServicePaths)),
		cron:     make([]wminit.Cron, 0, len(cfg.Cron)),
	}

	for _, s := range cfg.ServicePaths {
		retval.services = append(retval.services, wPath{path: s})
	}

	for _, c := range cfg.Cron {
		retval.cron = append(retval.cron, wCron{iCron: c})
	}

	return retval
}
