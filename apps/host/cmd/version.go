package cmd

import (
	"fmt"
	"gearbox/app"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of this CLI app.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(app.VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
