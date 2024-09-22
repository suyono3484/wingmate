package main

import (
	_ "embed"
	"os"

	"gitea.suyono.dev/suyono/wingmate"
	"gitea.suyono.dev/suyono/wingmate/config"
	wminit "gitea.suyono.dev/suyono/wingmate/init"
	"github.com/spf13/viper"
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
	config.SetVersion(version)
	config.ParseFlags()

	wingmate.Log().Info().Msgf("starting wingmate version %s", viper.GetString(config.WingmateVersion))

	if cfg, err = config.Read(); err != nil {
		wingmate.Log().Fatal().Err(err).Msg("failed to read config")
	}

	initCfg := convert(cfg)
	wminit.NewInit(initCfg).Start()
}
