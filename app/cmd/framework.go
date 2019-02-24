package cmd

import (
	"github.com/spf13/cobra"
)

var frameworkCmd = &cobra.Command{
	Use:   "framework",
	Short: "Manage Frameworks",
}

func init() {
	RootCmd.AddCommand(frameworkCmd)
	frameworkCmd.AddCommand(&cobra.Command{
		Use:   "load",
		Args:  cobra.ExactArgs(1),
		Short: "Load a framework",
		Run: func(cmd *cobra.Command, args []string) {
			//gearbox.Instance.LoadPlugins()
		},
	})
}
