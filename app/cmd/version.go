package cmd

import (
	"fmt"
	"gearbox/gearbox"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of this command line app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gearbox.VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
