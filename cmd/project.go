package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Gearbox projects.",
}

func init() {

	RootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Args:  cobra.ExactArgs(1),
		Short: "Initialize a Gearbox project.",
		Long:  "Initialize a Gearbox project by creating a default `project.json` file.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initialize Gearbox project goes here.")
		},
	})

	projectCmd.AddCommand(&cobra.Command{
		Use:   "enable",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Enable the current or specified Gearbox project.",
		Long: "Enable the current or specified Gearbox project will mark the project to run while Gearbox is up and running, " +
			"will enable any required services, and will start any of the required services if not already enabled.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Enable Gearbox project goes here...")
		},
	})

	projectCmd.AddCommand(&cobra.Command{
		Use:   "disable",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Disable the current or specified Gearbox project.",
		Long: "Disable the current or specified Gearbox project will mark the project to not run while Gearbox is up and " +
			"running, will disable any services required for this project not needed by any remaining enabled services, " +
			"and will stop any of those services it finds that need to be disabled.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Disable Gearbox service containers goes here...")
		},
	})
}
