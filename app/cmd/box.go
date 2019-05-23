package cmd

import (
	"fmt"
	"gearbox/box"
	"gearbox/gearbox"
	"gearbox/ssh"
	"github.com/gearboxworks/go-status/is"
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
	var boxArgs box.Args

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
			sts := gearbox.Instance.StartBox(boxArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
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
			sts := gearbox.Instance.StopBox(boxArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	boxCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Display the current status of the Gearbox VM.",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.PrintBoxStatus(boxArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
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
			sts := gearbox.Instance.RestartBox(boxArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	boxCmd.AddCommand(&cobra.Command{
		Use: "create",
		SuggestFor: []string{
			"new",
		},
		Short: "Create a new instance of the Gearbox VM.",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.CreateBox(boxArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	boxCmd.PersistentFlags().BoolVarP(&boxArgs.NoWait, "no-wait", "", false, "Don't wait for Box (VM) operation to complete.")
	boxCmd.PersistentFlags().IntVarP(&boxArgs.WaitRetries, "wait-delay", "", box.DefaultWaitRetries, "Box (VM) operation wait retries.")
	boxCmd.PersistentFlags().DurationVarP(&boxArgs.WaitDelay, "wait-retries", "", box.DefaultWaitDelay, "Box (VM) operation wait delay between retries.")
	boxCmd.PersistentFlags().StringVarP(&boxArgs.ConsoleHost, "console-host", "", box.DefaultConsoleHost, "Box (VM) console host name.")
	boxCmd.PersistentFlags().StringVarP(&boxArgs.ConsolePort, "console-port", "", box.DefaultConsolePort, "Box (VM) console port number.")
	boxCmd.PersistentFlags().BoolVarP(&boxArgs.ShowConsole, "show-console", "", box.DefaultShowConsole, "Show Box (VM) console output.")

	boxCmd.PersistentFlags().StringVarP(&boxArgs.SshUsername, "user", "u", ssh.DefaultUsername, "Alternate Gearbox SSH username.") // boxCmd.PersistentFlags().BoolP("no-wait", "w", false, "Don't wait for Box (VM) operation to complete.")
	boxCmd.PersistentFlags().StringVarP(&boxArgs.SshPassword, "password", "p", ssh.DefaultPassword, "Alternate Gearbox SSH password.")
	boxCmd.PersistentFlags().StringVarP(&boxArgs.SshPublicKey, "key-file", "k", ssh.DefaultKeyFile, "Gearbox SSH public key file.") // boxCmd.Flag("no-wait")
}
