package main

import (
	_ "embed"
	"os"

	"gitea.suyono.dev/suyono/wingmate"
	"gitea.suyono.dev/suyono/wingmate/config"
	wminit "gitea.suyono.dev/suyono/wingmate/init"
)

var (
	//go:embed version.txt
	version string
)

func main() {
	var (
		err error
		cfg *config.Config
	)

	_ = wingmate.NewLog(os.Stderr)
	if cfg, err = config.Read(); err != nil {
		wingmate.Log().Error().Msgf("failed to read config %#v", err)
	}

	initCfg := convert(cfg)
	wminit.NewInit(initCfg).Start()
}
