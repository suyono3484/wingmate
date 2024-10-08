package main

import (
	"log"
	"os/exec"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
)

const (
	EnvOneShotPath = "ONESHOT_PATH"
	OneShotPath    = "/usr/local/bin/wmoneshot"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
		t   *time.Ticker
	)
	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvOneShotPath)
	viper.SetDefault(EnvOneShotPath, OneShotPath)

	exePath := viper.GetString(EnvOneShotPath)

	t = time.NewTicker(time.Second * 5)
	for {
		cmd = exec.Command(exePath, "--", "wmdummy")
		if err = cmd.Run(); err != nil {
			log.Printf("failed to run %s: %+v\n", exePath, err)
		} else {
			log.Printf("%s executed\n", exePath)
		}

		<-t.C
	}
}
