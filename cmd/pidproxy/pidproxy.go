package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

const (
	pidFileFlag         = "pid-file"
	EnvStartSecs        = "STARTSECS"
	EnvDefaultStartSecs = 1
)

var (
	rootCmd = &cobra.Command{
		Use:  "wmpidproxy",
		RunE: pidProxy,
	}

	childArgs []string
)

func main() {
	var (
		i        int
		arg      string
		selfArgs []string
		found    bool
	)

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvStartSecs)
	viper.SetDefault(EnvStartSecs, EnvDefaultStartSecs)

	rootCmd.PersistentFlags().StringP(pidFileFlag, "p", "", "location of pid file")
	rootCmd.MarkFlagRequired(pidFileFlag)
	viper.BindPFlag(pidFileFlag, rootCmd.PersistentFlags().Lookup(pidFileFlag))

	found = false
	for i, arg = range os.Args {
		if arg == "--" {
			found = true
			if len(os.Args) <= i+1 {
				log.Println("invalid argument")
				os.Exit(1)
			}
			selfArgs = os.Args[1:i]
			childArgs = os.Args[i+1:]
			break
		}
	}
	if !found {
		log.Println("invalid argument")
		os.Exit(1)
	}

	if len(childArgs) == 0 {
		log.Println("invalid argument")
		os.Exit(1)
	}

	rootCmd.SetArgs(selfArgs)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func pidProxy(cmd *cobra.Command, args []string) error {
	pidfile := viper.GetString(pidFileFlag)
	log.Printf("%s %v", pidfile, childArgs)
	if len(childArgs) > 1 {
		go startProcess(childArgs[0], childArgs[1:]...)
	} else {
		go startProcess(childArgs[0])
	}
	initialWait := viper.GetInt(EnvStartSecs)
	time.Sleep(time.Second * time.Duration(initialWait))

	var (
		err error
		pid int
		sc  chan os.Signal
		t   *time.Timer
	)

	sc = make(chan os.Signal, 1)
	signal.Notify(sc, unix.SIGTERM)

	t = time.NewTimer(time.Second)

check:
	for {
		if pid, err = readPid(pidfile); err != nil {
			return err
		}

		if err = unix.Kill(pid, syscall.Signal(0)); err != nil {
			return err
		}

		select {
		case <-t.C:
		case <-sc:
			if pid, err = readPid(pidfile); err != nil {
				return err
			}

			if err = unix.Kill(pid, unix.SIGTERM); err != nil {
				return err
			}
			break check
		}
	}
	return nil
}

func readPid(pidFile string) (int, error) {
	var (
		file  *os.File
		err   error
		pid64 int64
	)

	if file, err = os.Open(pidFile); err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		if pid64, err = strconv.ParseInt(scanner.Text(), 10, 64); err != nil {
			return 0, err
		}
		return int(pid64), nil
	} else {
		return 0, errors.New("invalid scanner")
	}
}

func startProcess(arg0 string, args ...string) {
	if err := exec.Command(arg0, args...).Run(); err != nil {
		log.Println("exec:", err)
		return
	}
}
