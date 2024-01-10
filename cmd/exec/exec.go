package main

import (
	_ "embed"
	"errors"
	"fmt"
	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type execApp struct {
	childArgs []string
	err       error
}

const (
	setsidFlag  = "setsid"
	EnvSetsid   = "SETSID"
	userFlag    = "user"
	EnvUser     = "USER"
	versionFlag = "version"
)

var (

	//go:embed version.txt
	version string
)

func main() {
	var (
		selfArgs   []string
		childArgs  []string
		app        *execApp
		rootCmd    *cobra.Command
		versionCmd *cobra.Command
		err        error
	)

	app = &execApp{}

	rootCmd = &cobra.Command{
		Use:          "wmexec",
		SilenceUsage: true,
		RunE:         app.execCmd,
	}

	versionCmd = &cobra.Command{
		Use:  "version",
		RunE: app.versionCmd,
	}

	rootCmd.PersistentFlags().BoolP(setsidFlag, "s", false, "set to true to run setsid() before exec")
	viper.BindPFlag(EnvSetsid, rootCmd.PersistentFlags().Lookup(setsidFlag))

	rootCmd.PersistentFlags().StringP(userFlag, "u", "", "\"user:[group]\"")
	viper.BindPFlag(EnvUser, rootCmd.PersistentFlags().Lookup(userFlag))

	rootCmd.PersistentFlags().Bool(versionFlag, false, "print version")
	viper.BindPFlag(versionFlag, rootCmd.PersistentFlags().Lookup(versionFlag))

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvUser)
	viper.BindEnv(EnvSetsid)
	viper.SetDefault(EnvSetsid, false)
	viper.SetDefault(EnvUser, "")

	rootCmd.AddCommand(versionCmd)

	selfArgs, childArgs, err = argSplit()
	app.childArgs = childArgs
	app.err = err
	rootCmd.SetArgs(selfArgs)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func argSplit() ([]string, []string, error) {
	var (
		i         int
		arg       string
		selfArgs  []string
		childArgs []string
	)
	found := false
	for i, arg = range os.Args {
		if arg == "--" {
			found = true
			if i+1 == len(os.Args) {
				return nil, nil, errors.New("invalid argument")
			}

			if len(os.Args[i+1:]) == 0 {
				return nil, nil, errors.New("invalid argument")
			}

			selfArgs = os.Args[1:i]
			childArgs = os.Args[i+1:]
			break
		}

		if !found {
			return nil, nil, errors.New("invalid argument")
		}
	}
	return selfArgs, childArgs, nil
}

func (e *execApp) versionCmd(cmd *cobra.Command, args []string) error {
	e.printVersion()
	return nil
}

func (e *execApp) printVersion() {
	fmt.Print(version)
	os.Exit(0)
}

func (e *execApp) execCmd(cmd *cobra.Command, args []string) error {
	if viper.GetBool(versionFlag) {
		e.printVersion()
	}

	if e.err != nil {
		return e.err
	}

	if viper.GetBool(EnvSetsid) {
		_, _ = unix.Setsid()
	}

	var (
		uid  uint64
		gid  uint64
		err  error
		path string
	)

	ug := viper.GetString(EnvUser)
	if len(ug) > 0 {
		user, group, ok := strings.Cut(ug, ":")
		if ok {
			if gid, err = strconv.ParseUint(group, 10, 32); err != nil {
				if gid, err = getGid(group); err != nil {
					return fmt.Errorf("cgo getgid: %w", err)
				}
			}
			if err = unix.Setgid(int(gid)); err != nil {
				return fmt.Errorf("setgid: %w", err)
			}
		}

		uid, err = strconv.ParseUint(user, 10, 32)
		if err != nil {
			if uid, err = getUid(user); err != nil {
				return fmt.Errorf("cgo getuid: %w", err)
			}
		}
		if err = unix.Setuid(int(uid)); err != nil {
			return fmt.Errorf("setuid: %w", err)
		}

	}

	if path, err = exec.LookPath(e.childArgs[0]); err != nil {
		if !errors.Is(err, exec.ErrDot) {
			return fmt.Errorf("lookpath: %w", err)
		}
	}

	if err = unix.Exec(path, e.childArgs, os.Environ()); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
