package cmd

import (
	"fmt"
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/types"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Parent projects.",
}

func init() {

	RootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Args:  cobra.ExactArgs(1),
		Short: "Initialize a Parent project.",
		Long:  "Initialize a Parent project by creating a default `project.json` file.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initialize Parent project goes here.")
		},
	})

	projectCmd.AddCommand(&cobra.Command{
		Use:   "enable",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Enable the current or specified Parent project.",
		Long: "Enable the current or specified Parent project will mark the project to run while Parent is up and running, " +
			"will enable any required services, and will start any of the required services if not already enabled.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Enable Parent project goes here...")
		},
	})

	projectCmd.AddCommand(&cobra.Command{
		Use:   "disable",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Disable the current or specified Parent project.",
		Long: "Disable the current or specified Parent project will mark the project to not run while Parent is up and " +
			"running, will disable any services required for this project not needed by any remaining enabled services, " +
			"and will stop any of those services it finds that need to be disabled.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Disable Parent service containers goes here...")
		},
	})

	projectRootCmd := &cobra.Command{
		Use:   "root",
		Args:  cobra.ExactArgs(1),
		Short: "Manage project root directories on a host computer.",
	}

	projectCmd.AddCommand(projectRootCmd)

	config.ProjectRootAddCmd = &cobra.Command{
		Use:   "add",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Add a new directory on a host computer to the list of project roots in use.",
		Run: func(cmd *cobra.Command, args []string) {
			var dir types.AbsoluteDir
			var basedir types.Nickname
			switch len(args) {
			case 2:
				basedir = types.Nickname(args[1])
				fallthrough
			case 1:
				dir = types.AbsoluteDir(args[0])
			}
			if dir != "" {
				sts := gearbox.Instance.AddBasedir(dir, basedir)
				fmt.Printf(sts.Message())
			}
		},
	}

	projectRootCmd.AddCommand(config.ProjectRootAddCmd)

}
