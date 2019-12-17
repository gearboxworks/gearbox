package cmd

import (
	gopt "gearbox/global"
	"github.com/spf13/cobra"
)

var GlobalOptions *gopt.Options

var RootCmd = &cobra.Command{
	Use:   "gearbox",
	Short: "Manage your GearboxOS virtual machine.",
}

func init() {
	GlobalOptions = &gopt.Options{}
	pf := RootCmd.PersistentFlags()
	pf.BoolVarP(&GlobalOptions.IsDebug, "debug", "", false, "Debug mode")
	pf.BoolVarP(&GlobalOptions.NoCache, "no-cache", "", false, "Disable caching")
}
