package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use: "ssh",
	SuggestFor: []string{
		"login",
		"logon",
		"access",
	},
	Short: "Connect to the terminal of the running GearboxOS VM.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gearbox SSH code goes here...")
	},
}

func init() {
	RootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
