package main

import (
	"log"
	"math/rand"
	"os"
	"os/exec"

	"gitea.suyono.dev/suyono/wingmate"
	"gitea.suyono.dev/suyono/wingmate/cmd/cli"
	"github.com/spf13/viper"
)

const (
	EnvLog               = "LOG"
	EnvLogMessage        = "LOG_MESSAGE"
	EnvDefaultLogMessage = "oneshot executed"
	EnvInstanceNum       = "INSTANCE_NUM"
	EnvDefaultInstances  = 0
)

func main() {
	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvLog)
	viper.BindEnv(EnvLogMessage)
	viper.BindEnv(EnvInstanceNum)
	viper.SetDefault(EnvLogMessage, EnvDefaultLogMessage)
	viper.SetDefault(EnvInstanceNum, EnvDefaultInstances)

	_, childArgs, err := cli.SplitArgs(os.Args)
	if err != nil {
		log.Printf("splitargs: %+v", err)
		os.Exit(2)
	}

	logPath := viper.GetString(EnvLog)
	logMessage := viper.GetString(EnvLogMessage)
	log.Println("log path:", logPath)
	if logPath != "" {
		var (
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

	if len(childArgs) > 0 {
		StartInstances(childArgs[0], childArgs[1:]...)
	}
}

func StartInstances(exePath string, args ...string) {
	num := (rand.Uint32() % 16) + 16

	iNum := viper.GetInt(EnvInstanceNum)
	if iNum > 0 {
		num = uint32(iNum)
	}

	var (
		ctr uint32
		cmd *exec.Cmd
		err error
	)
	for ctr = 0; ctr < num; ctr++ {
		cmd = exec.Command(exePath, args...)
		if err = cmd.Start(); err != nil {
			log.Printf("failed to run %s: %+v\n", exePath, err)
		}
	}
}
