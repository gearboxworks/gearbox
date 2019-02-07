package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use: "stop",
	SuggestFor: []string{
		"poweroff",
		"shutdown",
		"down",
		"halt",
		"end",
		"off",
	},
	Short: "Stops the Gearbox VM if it is running.",
	Long: "The `gearbox stop` command contact VirtualBox and requests that it shutdown the GearboxOS " +
		"virtual machine that should be running within VirtualBox.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gearbox shutdown code goes here...")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
