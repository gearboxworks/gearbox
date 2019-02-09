package cmd

import (
	"gearbox/gearbox"
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Load the Gearbox admin console.",
	Run: func(cmd *cobra.Command, args []string) {

		admin := &gearbox.Admin{
			Host: gearbox.Host,
		}
		admin.Run()
	},
}

func init() {
	RootCmd.AddCommand(adminCmd)
}
