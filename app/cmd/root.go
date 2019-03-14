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
	pf := RootCmd.PersistentFlags()
	pf.BoolVarP(&GlobalOptions.IsDebug, "debug", "", false, "Debug mode")
	pf.BoolVarP(&GlobalOptions.NoCache, "no-cache", "", false, "Disable caching")
}
