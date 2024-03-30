package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ParseFlags() {
	pflag.CountP(VersionFlag, "v", "check version")
	pflag.String(WMPidProxyPathFlag, "", "wmpidproxy path")
	pflag.String(WMExecPathFlag, "", "wmexec path")
	pflag.StringP(ConfigPathFlag, "c", "", "config path")

	pflag.Parse()

	_ = viper.BindPFlag(VersionCheckKey, pflag.CommandLine.Lookup(VersionFlag))
	_ = viper.BindPFlag(ConfigPath, pflag.CommandLine.Lookup(ConfigPathFlag))
	_ = viper.BindPFlag(PidProxyPathConfig, pflag.CommandLine.Lookup(WMPidProxyPathFlag))
	_ = viper.BindPFlag(ExecPathConfig, pflag.CommandLine.Lookup(WMExecPathFlag))

	if viper.GetInt(VersionCheckKey) > 0 {
		fmt.Println(viper.GetString(WingmateVersion))
		os.Exit(0)
	}
}
