// +build !vm
package cmd

import (
	"gearbox"
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Load the Gearbox admin console.",
	Run: func(cmd *cobra.Command, args []string) {
		gb := gearbox.Instance
		gb.HostConnector = gearbox.Host
		gb.Admin()
	},
}

func init() {
	RootCmd.AddCommand(adminCmd)
}
