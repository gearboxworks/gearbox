// +build !vm

package cmd

import (
	"gearbox/gearbox"
	"github.com/spf13/cobra"
)

type adminArgs struct {
	Viewer string
}

var AdminArgs = adminArgs{}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Load the Gearbox admin console.",
	Run: func(cmd *cobra.Command, args []string) {
		gb := gearbox.Instance
		gb.Admin(gearbox.ViewerType(AdminArgs.Viewer))
	},
}

func init() {
	RootCmd.AddCommand(adminCmd)

	adminCmd.Flags().StringVarP(
		&AdminArgs.Viewer,
		"viewer",
		"",
		string(gearbox.DefaultViewer),
		"Web viewer to use: 'lorca' or 'webview'",
	)

}
