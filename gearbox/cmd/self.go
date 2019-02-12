package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage this command line app",
}

func init() {
	RootCmd.AddCommand(selfCmd)
	selfCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update this command line app to the latest compatible version",
		Long:  "Update this `gearbox` command line app to the latest version compatible with the current version of Gearbox OS.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Update this `gearbox` CLI app  goes here.")
		},
	})
	selfCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Run tests for GearBox itself (We have many yet to write. Wanna help?)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run tests for GearBox itself goes here.")
		},
	})
	selfCmd.AddCommand(&cobra.Command{
		Use:   "checklist",
		Short: "Our launch checklist for new versions of GearBox (we eat our own dog food)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Our launch checklist goes here.")
		},
	})
}
