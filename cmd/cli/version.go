package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Version string

const versionFlag = "version"

func (v Version) Print() {
	fmt.Print(v)
	os.Exit(0)
}

func (v Version) Cmd(cmd *cobra.Command) {
	cmd.AddCommand(&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			v.Print()
			return nil
		},
	})
}

func (v Version) Flag(cmd *cobra.Command) {
	viper.SetDefault(versionFlag, false)
	cmd.PersistentFlags().Bool(versionFlag, false, "print version")
	_ = viper.BindPFlag(versionFlag, cmd.PersistentFlags().Lookup(versionFlag))
}

func (v Version) FlagHook() {
	if viper.GetBool(versionFlag) {
		v.Print()
	}
}
