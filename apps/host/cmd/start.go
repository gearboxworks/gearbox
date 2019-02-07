package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	SuggestFor: []string{
		"startup",
		"begin",
		"run",
		"up",
		"on",
	},
	Short: "Starts up the Gearbox VM running Gearbox OS.",
	Long: "The `gearbox start` command is used to 'start' gearbox. It runs checks to ensure everything " +
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
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
