package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"gitea.suyono.dev/suyono/wingmate"
	"gitea.suyono.dev/suyono/wingmate/cmd/cli"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DummyPath    = "/usr/local/bin/wmdummy"
	EnvDummyPath = "DUMMY_PATH"
	NoWaitFlag   = "no-wait"
)

func main() {
	var (
		stdout    io.ReadCloser
		stderr    io.ReadCloser
		wg        *sync.WaitGroup
		err       error
		exePath   string
		selfArgs  []string
		childArgs []string
		flagSet   *pflag.FlagSet
		noWait    bool
		cmd       *exec.Cmd
	)
	if selfArgs, childArgs, err = cli.SplitArgs(os.Args); err == nil {
		flagSet = pflag.NewFlagSet(selfArgs[0], pflag.ExitOnError)
		flagSet.Bool(NoWaitFlag, false, "do not wait for the child process")
		if err = flagSet.Parse(selfArgs[1:]); err != nil {
			log.Printf("invalid argument: %+v", err)
			return
		}
	} else {
		flagSet = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
		flagSet.Bool(NoWaitFlag, false, "do not wait for the child process")
		if err = flagSet.Parse(selfArgs[1:]); err != nil {
			log.Printf("invalid argument: %+v", err)
			return
		}
	}
	viper.BindPFlag(NoWaitFlag, flagSet.Lookup(NoWaitFlag))
	noWait = viper.GetBool(NoWaitFlag)

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvDummyPath)
	viper.SetDefault(EnvDummyPath, DummyPath)

	exePath = viper.GetString(EnvDummyPath)

	if len(childArgs) > 0 {
		cmd = exec.Command(childArgs[0], childArgs[1:]...)
	} else {
		cmd = exec.Command(exePath)
	}

	if !noWait {
		if stdout, err = cmd.StdoutPipe(); err != nil {
			log.Panic(err)
		}

		if stderr, err = cmd.StderrPipe(); err != nil {
			log.Panic(err)
		}

		wg = &sync.WaitGroup{}
		wg.Add(2)
		go pulley(wg, stdout, "stdout")
		go pulley(wg, stderr, "stderr")
	}

	if err = cmd.Start(); err != nil {
		log.Panic(err)
	}

	if !noWait {
		wg.Wait()

		if err = cmd.Wait(); err != nil {
			log.Printf("got error when Waiting for child process: %#v\n", err)
		}
	}
}

func pulley(wg *sync.WaitGroup, src io.ReadCloser, srcName string) {
	defer wg.Done()

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		log.Printf("coming out from %s: %s\n", srcName, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("got error whean reading from %s: %#v\n", srcName, err)
	}

	log.Printf("closing %s...\n", srcName)
	_ = src.Close()
}
