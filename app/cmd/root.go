package cmd

import (
	"gearbox"
	"github.com/spf13/cobra"
)

var GlobalOptions *gearbox.GlobalOptions

var RootCmd = &cobra.Command{
	Use:   "gearbox",
	Short: "Manage your GearboxOS virtual machine.",
}

func init() {
	GlobalOptions = &gearbox.GlobalOptions{}
	RootCmd.PersistentFlags().BoolVarP(&GlobalOptions.IsDebug, "debug", "", false, "Debug mode")
}
