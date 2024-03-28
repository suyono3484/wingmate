package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/sys/unix"
)

const (
	logPathFlag = "log-path"
	pidFileFlag = "pid-file"
)

func main() {
	var (
		logPath     string
		pidFilePath string
		name        string
		pause       uint
		lf          *os.File
		err         error
	)

	pflag.StringVarP(&logPath, logPathFlag, "l", "/var/log/wmbg.log", "log file path")
	pflag.StringVarP(&pidFilePath, pidFileFlag, "p", "/var/run/wmbg.pid", "pid file path")
	pflag.StringVar(&name, "name", "no-name", "process name")
	pflag.UintVar(&pause, "pause", 5, "pause interval")
	pflag.Parse()

	if lf, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		os.Exit(2)
	}
	defer func() {
		_ = lf.Close()
	}()

	log.SetOutput(lf)
	log.Printf("starting process %s with pause interval %d", name, pause)
	if err = writePid(pidFilePath); err != nil {
		log.Printf("failed to write pid file: %+v", err)
	}
	time.Sleep(time.Duration(pause) * time.Second)
	log.Printf("process %s finished", name)
}

func writePid(path string) error {
	var (
		err error
		pf  *os.File
	)

	if pf, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("opening pid file for write: %w", err)
	}
	defer func() {
		_ = pf.Close()
	}()

	if _, err = fmt.Fprintf(pf, "%d", unix.Getpid()); err != nil {
		return fmt.Errorf("writing pid to the pid file: %w", err)
	}
	return nil
}
