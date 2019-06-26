package cmd

import (
	"fmt"
	gopt "gearbox/global"
	"github.com/spf13/cobra"
	"os"
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
	pf.BoolVarP(&GlobalOptions.NoDownloadGears, "no-download-gears", "", false, "Disable downloading gears.json")
	err := pf.Parse(os.Args)
	if err != nil {
		panic(fmt.Sprintf("unable to set CLI flags: %s", err))
	}
}
