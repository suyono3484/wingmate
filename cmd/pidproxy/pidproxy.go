package main

import (
	"bufio"
	"errors"
	"gitea.suyono.dev/suyono/wingmate/cmd/cli"
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

	_ "embed"
)

type pidProxyApp struct {
	childArgs []string
	err       error
	version   cli.Version
}

const (
	pidFileFlag         = "pid-file"
	EnvStartSecs        = "STARTSECS"
	EnvDefaultStartSecs = 1
)

var (

	//go:embed version.txt
	version string
)

func main() {
	var (
		selfArgs  []string
		childArgs []string
		err       error
		app       *pidProxyApp
		rootCmd   *cobra.Command
	)

	app = &pidProxyApp{
		version: cli.Version(version),
	}

	rootCmd = &cobra.Command{
		Use:          "wmpidproxy",
		SilenceUsage: true,
		RunE:         app.pidProxy,
	}

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	_ = viper.BindEnv(EnvStartSecs)
	viper.SetDefault(EnvStartSecs, EnvDefaultStartSecs)

	rootCmd.PersistentFlags().StringP(pidFileFlag, "p", "", "location of pid file")
	_ = rootCmd.MarkFlagRequired(pidFileFlag)
	_ = viper.BindPFlag(pidFileFlag, rootCmd.PersistentFlags().Lookup(pidFileFlag))

	app.version.Flag(rootCmd)
	app.version.Cmd(rootCmd)

	if selfArgs, childArgs, err = cli.SplitArgs(os.Args); err != nil {
		selfArgs = os.Args
	}
	app.childArgs = childArgs
	app.err = err

	rootCmd.SetArgs(selfArgs[1:])
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (p *pidProxyApp) pidProxy(cmd *cobra.Command, args []string) error {
	p.version.FlagHook()

	pidfile := viper.GetString(pidFileFlag)
	log.Printf("%s %v", pidfile, p.childArgs)
	if len(p.childArgs) > 1 {
		go p.startProcess(p.childArgs[0], p.childArgs[1:]...)
	} else {
		go p.startProcess(p.childArgs[0])
	}
	initialWait := viper.GetInt(EnvStartSecs)
	time.Sleep(time.Second * time.Duration(initialWait))

	var (
		err error
		pid int
		sc  chan os.Signal
		t   *time.Ticker
	)

	sc = make(chan os.Signal, 1)
	signal.Notify(sc, unix.SIGTERM)

	t = time.NewTicker(time.Second)

check:
	for {
		if pid, err = p.readPid(pidfile); err != nil {
			return err
		}

		if err = unix.Kill(pid, syscall.Signal(0)); err != nil {
			if !errors.Is(err, unix.ESRCH) {
				return err
			}
			break check
		}

		select {
		case <-t.C:
		case <-sc:
			if pid, err = p.readPid(pidfile); err != nil {
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

func (p *pidProxyApp) readPid(pidFile string) (int, error) {
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

func (p *pidProxyApp) startProcess(arg0 string, args ...string) {
	if err := exec.Command(arg0, args...).Run(); err != nil {
		log.Println("exec:", err)
		return
	}
}
