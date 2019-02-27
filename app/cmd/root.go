package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gearbox",
	Short: "Manage your GearboxOS virtual machine.",
}
