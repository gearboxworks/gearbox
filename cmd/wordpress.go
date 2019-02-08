package cmd

import (
	"github.com/spf13/cobra"
)

var wordPressCmd = &cobra.Command{
	Use: "wordpress",
	Aliases: []string{
		"wp",
	},
	Short: "Manage Gearboxes' WordPress-specific options.",
}

func init() {
	RootCmd.AddCommand(wordPressCmd)
}
