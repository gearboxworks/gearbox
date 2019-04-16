package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Manage Parent service containers.",
}

func init() {
	RootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("List Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Args:  cobra.ExactArgs(1),
		Short: "Start Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Start Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Args:  cobra.ExactArgs(1),
		Short: "Stop Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stop Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use: "enable",
		SuggestFor: []string{
			"activate",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Activate Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Activate Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use: "disable",
		SuggestFor: []string{
			"deactivate",
		},
		Args:  cobra.ExactArgs(1),
		Short: "Disable Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("disable Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Args:  cobra.ExactArgs(1),
		Short: "Install a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Install a specific Parent service container goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Args:  cobra.ExactArgs(1),
		Short: "Uninstall a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Uninstall Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "upgrade",
		Args:  cobra.ExactArgs(1),
		Short: "Upgrade a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Upgrade Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "shutdown",
		Args:  cobra.NoArgs,
		Short: "Shutdown all Parent service containers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Shutdown Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "inspect",
		Args:  cobra.ExactArgs(1),
		Short: "Inspect details about a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "log",
		Args:  cobra.ExactArgs(1),
		Short: "Display logs for a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Parent service containers goes here...")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: "Inspect details about a specific Parent service container.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("THIS COMMAND SHOULD NOT BE NEEDED. " +
				"IT SHOULD UPDATED IN THE BACKGROUND.")
		},
	})
	serviceCmd.AddCommand(&cobra.Command{
		Use:   "pull",
		Args:  cobra.ExactArgs(1),
		Short: "Pull a specific Parent service container without installing.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display log for Parent service containers goes here...")
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
		Short: "Display the current status of the running Parent services.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display the current status of Parent goes here.")
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
		Long: "Stops all project-enabled service containers for Parent and then starts them back up again." +
			"\n" +
			"\nThis is equivalent to running `gearbox services stop` and then `gearbox services start`.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Parent services restart code goes here...")
		},
	})

}
