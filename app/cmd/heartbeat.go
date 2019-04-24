package cmd

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/heartbeat"
	"gearbox/ssh"
	"gearbox/status/is"
	"github.com/spf13/cobra"
)

var heartbeatCmd = &cobra.Command{
	Use: "heartbeat",
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
	Short: "Manage the Gearbox Heartbeat",
}

func init() {
	var heartbeatArgs heartbeat.Args

	RootCmd.AddCommand(heartbeatCmd)

	heartbeatCmd.AddCommand(&cobra.Command{
		Use: "daemon",
		SuggestFor: []string{
			"systray",
		},
		Short: "Gearbox Heartbeat daemon",
		Long: "The `gearbox heartbeat start` command is used to run the Heartbeat daemon. " +
			"This maintains low-level communications with important Gearbox applications and tools. " +
			"It also provides a user systray for control of Gearbox. ",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.HeartbeatDaemon(heartbeatArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	heartbeatCmd.AddCommand(&cobra.Command{
		Use: "start",
		SuggestFor: []string{
			"startup",
			"begin",
			"run",
			"up",
			"on",
		},
		Short: "Starts up the Gearbox Heartbeat daemon",
		Long: "The `gearbox heartbeat start` command is used to run the Heartbeat daemon. " +
			"This maintains low-level communications with important Gearbox applications and tools. " +
			"It also provides a user systray for control of Gearbox. ",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.StartHeartbeat(heartbeatArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	heartbeatCmd.AddCommand(&cobra.Command{
		Use: "stop",
		SuggestFor: []string{
			"poweroff",
			"shutdown",
			"down",
			"halt",
			"end",
			"off",
		},
		Short: "Stops the Gearbox Heartbeat if it is running",
		Long: "The `gearbox heartbeat stop` command.",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.StopHeartbeat(heartbeatArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	heartbeatCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Display the current status of the Gearbox Heartbeat.",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.PrintHeartbeatStatus(heartbeatArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	heartbeatCmd.AddCommand(&cobra.Command{
		Use: "restart",
		SuggestFor: []string{
			"recycle",
			"renew",
			"refresh",
		},
		Short: "Stops the Gearbox Heartbeat and then starts it back up again",
		Long: "Stops the Gearbox Heartbeat and then starts it back up again." +
			"\n" +
			"\nThis is equivalent to running `gearbox heartbeat stop` and then `gearbox heartbeat start`.",
		Run: func(cmd *cobra.Command, args []string) {
			sts := gearbox.Instance.RestartHeartbeat(heartbeatArgs)
			if is.Error(sts) {
				fmt.Println(sts.Message())
			}
		},
	})

	heartbeatCmd.PersistentFlags().BoolVarP(&heartbeatArgs.NoWait, "no-wait", "", false, "Don't wait for Box (VM) operation to complete.")
	heartbeatCmd.PersistentFlags().IntVarP(&heartbeatArgs.WaitRetries, "wait-delay", "", heartbeat.DefaultWaitRetries, "Box (VM) operation wait retries.")
	heartbeatCmd.PersistentFlags().DurationVarP(&heartbeatArgs.WaitDelay, "wait-retries", "", heartbeat.DefaultWaitDelay, "Box (VM) operation wait delay between retries.")
	heartbeatCmd.PersistentFlags().StringVarP(&heartbeatArgs.ConsoleHost, "console-host", "", heartbeat.DefaultConsoleHost, "Box (VM) console host name.")
	heartbeatCmd.PersistentFlags().StringVarP(&heartbeatArgs.ConsolePort, "console-port", "", heartbeat.DefaultConsolePort, "Box (VM) console port number.")
	heartbeatCmd.PersistentFlags().BoolVarP(&heartbeatArgs.ShowConsole, "show-console", "", heartbeat.DefaultShowConsole, "Show Box (VM) console output.")

	// Mike will not like this bit.
	heartbeatCmd.PersistentFlags().StringVarP(&heartbeatArgs.SshUsername, "user", "u", ssh.DefaultUsername, "Alternate Gearheartbeat SSH username.") // heartbeatCmd.PersistentFlags().BoolP("no-wait", "w", false, "Don't wait for Box (VM) operation to complete.")
	heartbeatCmd.PersistentFlags().StringVarP(&heartbeatArgs.SshPassword, "password", "p", ssh.DefaultPassword, "Alternate Gearheartbeat SSH password.")
	heartbeatCmd.PersistentFlags().StringVarP(&heartbeatArgs.SshPublicKey, "key-file", "k", ssh.DefaultKeyFile, "Gearbox SSH public key file.") // heartbeatCmd.Flag("no-wait")
}
