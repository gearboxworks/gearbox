package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
	Use: "vm",
	SuggestFor: []string{
		"virtual-machine",
		"virtualmachine",
	},
	Short: "Manage the Gearbox virtual machine",
}

func init() {
	RootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(&cobra.Command{
		Use: "start",
		SuggestFor: []string{
			"startup",
			"begin",
			"run",
			"up",
			"on",
		},
		Short: "Starts up the Gearbox VM running Gearbox OS",
		Long: "The `gearbox vm start` command is used to load Gearbox VM running GearboxOS. " +
			"It runs checks to ensure everything " +
			"required to successfully run Gearbox is installed and configured correctly, and if not it prompts " +
			"the user for permission to initiate and/or fix the configuration. In the case of the required dependency of " +
			"VirtualBox if it does not find it or it finds an incompatible version it will ask the user to install a " +
			"compatible version (but in the future we plan for this CLI app to install and configure VirtualBox for " +
			"the user. Once it ensures the configuration is correct it then starts VirtualBox, if needed, and requests" +
			"that VirtualBox start the ISO containing GearboxOS. Finally, once Gearbox OS is running and ready to serve " +
			"web requests `gearbox start` will tell the user that Gearbox is ready for use.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gearbox startup code goes here...")
		},
	})
	vmCmd.AddCommand(&cobra.Command{
		Use: "stop",
		SuggestFor: []string{
			"poweroff",
			"shutdown",
			"down",
			"halt",
			"end",
			"off",
		},
		Short: "Stops the Gearbox VM if it is running",
		Long: "The `gearbox stop` command contact VirtualBox and requests that it shutdown the GearboxOS " +
			"virtual machine that should be running within VirtualBox.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gearbox shutdown code goes here...")
		},
	})
	vmCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Display the current status of the Gearbox VM.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Display the current status of Gearbox VM goes here.")
		},
	})
	vmCmd.AddCommand(&cobra.Command{
		Use: "restart",
		SuggestFor: []string{
			"recycle",
			"renew",
			"refresh",
		},
		Short: "Stops the Gearbox virtual machine and then starts it back up again",
		Long: "Stops the Gearbox virtual machine and then starts it back up again." +
			"\n" +
			"\nThis is equivalent to running `gearbox vm stop` and then `gearbox vm start`.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gearbox vm restart code goes here...")
		},
	})
}
