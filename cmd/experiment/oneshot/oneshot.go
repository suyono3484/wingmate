package main

import (
	"log"
	"math/rand"
	"os"
	"os/exec"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
)

const (
	// DummyPath = "/workspaces/wingmate/cmd/experiment/dummy/dummy"
	DummyPath            = "/usr/local/bin/wmdummy"
	EnvDummyPath         = "DUMMY_PATH"
	EnvPrefix            = "WINGMATE"
	EnvLog               = "LOG"
	EnvLogMessage        = "LOG_MESSAGE"
	EnvDefaultLogMessage = "oneshot executed"
)

func main() {
	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv(EnvDummyPath)
	viper.BindEnv(EnvLog)
	viper.BindEnv(EnvLogMessage)
	viper.SetDefault(EnvDummyPath, DummyPath)
	viper.SetDefault(EnvLogMessage, EnvDefaultLogMessage)

	exePath := viper.GetString(EnvDummyPath)

	logPath := viper.GetString(EnvLog)
	logMessage := viper.GetString(EnvLogMessage)
	log.Println("log path:", logPath)
	if logPath != "" {
		var (
			err  error
			file *os.File
		)

		if file, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666); err == nil {
			defer func() {
				_ = file.Close()
			}()

			if err = wingmate.NewLog(file); err == nil {
				wingmate.Log().Info().Msg(logMessage)
			}
		}
	}

	StartRandomInstances(exePath)
}

func StartRandomInstances(exePath string) {
	num := (rand.Uint32() % 16) + 16

	var (
		ctr uint32
		cmd *exec.Cmd
		err error
	)
	for ctr = 0; ctr < num; ctr++ {
		cmd = exec.Command(exePath)
		if err = cmd.Start(); err != nil {
			log.Printf("failed to run %s: %+v\n", exePath, err)
		}
	}
}
