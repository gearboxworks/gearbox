package cmd

import (
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Parent project to a web host.",
	Long:  "This is a planned feature of Parent. The more people adopt Parent, the sooner we can implement it",
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
