package cmd

import (
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Gearbox project to a web host.",
	Long:  "This is a planned feature of Gearbox. The more people adopt Gearbox, the sooner we can implement it",
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
