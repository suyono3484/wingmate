package config

import (
	"fmt"

	"gitea.suyono.dev/suyono/wingmate/cmd/cli"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ParseFlags() {
	version := cli.Version(fmt.Sprintln(viper.GetString(WingmateVersion)))
	version.FlagSet(pflag.CommandLine)

	pflag.String(WMPidProxyPathFlag, "", "wmpidproxy path")
	pflag.String(WMExecPathFlag, "", "wmexec path")
	pflag.StringP(PathConfigFlag, "c", "", "config path")

	pflag.Parse()

	_ = viper.BindPFlag(PathConfig, pflag.CommandLine.Lookup(PathConfigFlag))
	_ = viper.BindPFlag(PidProxyPathConfig, pflag.CommandLine.Lookup(WMPidProxyPathFlag))
	_ = viper.BindPFlag(ExecPathConfig, pflag.CommandLine.Lookup(WMExecPathFlag))

	version.FlagHook()
}
