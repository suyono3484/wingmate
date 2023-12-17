package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

const (
	setsidFlag = "setsid"
	EnvSetsid  = "SETSID"
	userFlag   = "user"
	EnvUser    = "USER"
)

var (
	rootCmd = &cobra.Command{
		Use:  "wmexec",
		RunE: execCmd,
	}

	childArgs []string
)

func main() {
	var (
		found    bool
		i        int
		arg      string
		selfArgs []string
	)

	rootCmd.PersistentFlags().BoolP(setsidFlag, "s", false, "set to true to run setsid() before exec")
	viper.BindPFlag(EnvSetsid, rootCmd.PersistentFlags().Lookup(setsidFlag))

	rootCmd.PersistentFlags().StringP(userFlag, "u", "", "\"user:[group]\"")
	viper.BindPFlag(EnvUser, rootCmd.PersistentFlags().Lookup(userFlag))

	viper.SetEnvPrefix(wingmate.EnvPrefix)
	viper.BindEnv(EnvUser)
	viper.BindEnv(EnvSetsid)
	viper.SetDefault(EnvSetsid, false)
	viper.SetDefault(EnvUser, "")

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

func execCmd(cmd *cobra.Command, args []string) error {
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
		uid, err = strconv.ParseUint(user, 10, 32)
		if err != nil {
			if uid, err = getUid(user); err != nil {
				return err
			}
		}
		if err = unix.Setuid(int(uid)); err != nil {
			return err
		}

		if ok {
			if gid, err = strconv.ParseUint(group, 10, 32); err != nil {
				if gid, err = getGid(group); err != nil {
					return err
				}
			}
			if err = unix.Setgid(int(gid)); err != nil {
				return err
			}
		}
	}

	if path, err = exec.LookPath(childArgs[0]); err != nil {
		if !errors.Is(err, exec.ErrDot) {
			return err
		}
	}

	if err = unix.Exec(path, childArgs, os.Environ()); err != nil {
		return err
	}

	return nil
}
