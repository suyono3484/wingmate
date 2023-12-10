package main

import (
	"log"
	"os"
	"os/exec"
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

	if len(os.Args) <= 2 {
		log.Println("invalid argument")
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringP(pidFileFlag, "p", "", "location of pid file")
	rootCmd.MarkFlagRequired(pidFileFlag)
	viper.BindPFlag(pidFileFlag, rootCmd.PersistentFlags().Lookup(pidFileFlag))

	found = false
	for i, arg = range os.Args {
		if arg == "--" {
			found = true
			selfArgs = os.Args[1:i]
			if len(os.Args) <= i+1 {
				log.Println("invalid argument")
				os.Exit(1)
			}
			childArgs = os.Args[i+1:]
			break
		}
	}
	if !found {
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
	)
	for {
		if pid, err = readPid(pidfile); err != nil {
			return err
		}

		if err = unix.Kill(pid, syscall.Signal(0)); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
}

func readPid(pidFile string) (int, error) {
	var (
		file  *os.File
		err   error
		buf   []byte
		n     int
		pid64 int64
	)

	if file, err = os.Open(pidFile); err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	buf = make([]byte, 1024)
	n, err = file.Read(buf)
	if err != nil {
		return 0, err
	}

	pid64, err = strconv.ParseInt(string(buf[:n]), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(pid64), nil
}

func startProcess(arg0 string, args ...string) {
	if err := exec.Command(arg0, args...).Run(); err != nil {
		log.Println("exec:", err)
		return
	}
}
