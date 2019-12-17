package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Manage Gearbox service containers.",
}

func init() {
	RootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Args:  cobra.ExactArgs(1),
		Short: "Start Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Start Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Args:  cobra.ExactArgs(1),
		Short: "Stop Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stop Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use: "enable",
		SuggestFor: []string{
			"activate",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Activate Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Activate Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use: "disable",
		SuggestFor: []string{
			"deactivate",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Disable Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("disable Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Args:  cobra.ExactArgs(1),
		Short: "Install a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Install a specific Gearbox service container goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Args:  cobra.ExactArgs(1),
		Short: "Uninstall a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Uninstall Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "upgrade",
		Args:  cobra.ExactArgs(1),
		Short: "Upgrade a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Upgrade Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "shutdown",
		Args:  cobra.NoArgs,
		Short: "Shutdown all Gearbox service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Shutdown Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "inspect",
		Args:  cobra.ExactArgs(1),
		Short: "Inspect details about a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "log",
		Args:  cobra.ExactArgs(1),
		Short: "Display logs for a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: "Inspect details about a specific Gearbox service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("THIS COMMAND SHOULD NOT BE NEEDED. " +
				"IT SHOULD UPDATED IN THE BACKGROUND.")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "pull",
		Args:  cobra.ExactArgs(1),
		Short: "Pull a specific Gearbox service container without installing.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Gearbox service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "clean",
		Args:  cobra.NoArgs,
		Short: "LET'S DISCUSS IF THIS IS NEEDED.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("LET'S DISCUSS IF THIS IS NEEDED AND IF SO WHAT IT DOES")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "available",
		Args:  cobra.NoArgs,
		Short: "MAYBE THIS SHOULD BE A SWITCH ON `gearbox list?`",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("MAYBE `gearbox available` SHOULD BE A SWITCH ON `gearbox list?`")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Display the current status of the running Gearbox services.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display the current status of Gearbox goes here.")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use: "restart",
		SuggestFor: []string{
			"recycle",
			"renew",
			"refresh",
		},
		Short: "Stops all project-enabled service containers and then starts them back up again.",
		Long: "Stops all project-enabled service containers for Gearbox and then starts them back up again." +
			"\n" +
			"\nThis is equivalent to running `gearbox services stop` and then `gearbox services start`.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gearbox services restart code goes here...")
		},
	})

}
