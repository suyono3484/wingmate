package main

import (
	_ "embed"
	"errors"
	"fmt"
	"gitea.suyono.dev/suyono/wingmate"
	"gitea.suyono.dev/suyono/wingmate/cmd/cli"
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
	version   cli.Version
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
		selfArgs  []string
		childArgs []string
		app       *execApp
		rootCmd   *cobra.Command
		err       error
	)

	app = &execApp{
		version: cli.Version(version),
	}

	rootCmd = &cobra.Command{
		Use:          "wmexec",
		SilenceUsage: true,
		RunE:         app.execCmd,
	}

	rootCmd.PersistentFlags().BoolP(setsidFlag, "s", false, "set to true to run setsid() before exec")
	viper.BindPFlag(EnvSetsid, rootCmd.PersistentFlags().Lookup(setsidFlag))

	rootCmd.PersistentFlags().StringP(userFlag, "u", "", "\"user:[group]\"")
	viper.BindPFlag(EnvUser, rootCmd.PersistentFlags().Lookup(userFlag))

	app.version.Flag(rootCmd)

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvUser)
	viper.BindEnv(EnvSetsid)
	viper.SetDefault(EnvSetsid, false)
	viper.SetDefault(EnvUser, "")

	app.version.Cmd(rootCmd)

	selfArgs, childArgs, err = cli.SplitArgs()
	app.childArgs = childArgs
	app.err = err

	rootCmd.SetArgs(selfArgs)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (e *execApp) execCmd(cmd *cobra.Command, args []string) error {
	e.version.FlagHook()

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
