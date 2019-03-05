package cmd

import (
	"gearbox"
	"github.com/spf13/cobra"
)

var boxCmd = &cobra.Command{
	Use: "box",
	SuggestFor: []string{
		"vm",
		"engine",
		"virtualbox",
		"vmware",
		"parallels",
		"machine",
		"virtual-machine",
		"virtualmachine",
	},
	Short: "Manage the Gearbox virtual machine",
}

func init() {
	var vmArgs gearbox.BoxArgs

	RootCmd.AddCommand(boxCmd)
	boxCmd.AddCommand(&cobra.Command{
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
			err := gearbox.Instance.StartBox(vmArgs)
			if err != nil {
				panic(err)
			}
		},
	})

	boxCmd.AddCommand(&cobra.Command{
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
			err := gearbox.Instance.StopBox(vmArgs)
			if err != nil {
				panic(err)
			}
		},
	})

	boxCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Display the current status of the Gearbox VM.",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := gearbox.Instance.GetBoxStatus(vmArgs)
			if err != nil {
				panic(err)
			}
		},
	})

	boxCmd.AddCommand(&cobra.Command{
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
			err := gearbox.Instance.RestartBox(vmArgs)
			if err != nil {
				panic(err)
			}
		},
	})

	boxCmd.PersistentFlags().BoolVarP(&vmArgs.NoWait, "no-wait", "", false, "Don't wait for VM operation to complete.")
	boxCmd.PersistentFlags().IntVarP(&vmArgs.WaitRetries, "wait-delay", "", gearbox.BoxDefaultWaitRetries, "VM operation wait retries.")
	boxCmd.PersistentFlags().DurationVarP(&vmArgs.WaitDelay, "wait-retries", "", gearbox.BoxDefaultWaitDelay, "VM operation wait delay between retries.")
	boxCmd.PersistentFlags().StringVarP(&vmArgs.ConsoleHost, "console-host", "", gearbox.BoxDefaultConsoleHost, "VM console host name.")
	boxCmd.PersistentFlags().StringVarP(&vmArgs.ConsolePort, "console-port", "", gearbox.BoxDefaultConsolePort, "VM console port number.")
	boxCmd.PersistentFlags().BoolVarP(&vmArgs.ShowConsole, "show-console", "", gearbox.BoxDefaultShowConsole, "Show VM console output.")
	boxCmd.PersistentFlags().StringVarP(&vmArgs.Name, "name", "", gearbox.BoxDefaultName, "Gearbox VM name.")

	// boxCmd.PersistentFlags().BoolP("no-wait", "w", false, "Don't wait for VM operation to complete.")

	// boxCmd.Flag("no-wait")
}
