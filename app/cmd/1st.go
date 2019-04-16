package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var firstCmd = &cobra.Command{
	Use:   "1st",
	Short: "See what commands you need to run to get started",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nIt really is super simple to get started using Parent. Just type " +
			"\nthe following, press <Enter>, and then follow any instructions:" +
			"\n" +
			"\n\tgearbox vm start")
	},
}

func init() {
	RootCmd.AddCommand(firstCmd)
}
