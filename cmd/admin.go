package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zserge/webview"
	"net/url"
)

const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Gearbox Admin Console</title>
</head>
<body>
<h1>Gearbox FTW!</h1>

</body>
</html>`

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Load the Gearbox admin console.",
	Run: func(cmd *cobra.Command, args []string) {
		//
		// See https://github.com/zserge/webview
		//
		w := webview.New(webview.Settings{
			Title:     "Gearbox Admin Console",
			Height:    600,
			Width:     800,
			Resizable: true,
			URL:       `data:text/html,` + url.PathEscape(html),
		})
		w.Run()
	},
}

func init() {
	RootCmd.AddCommand(adminCmd)
}
