package main

import (
	"fmt"
	"gearbox"
	"gearbox/app/cmd"
	"gearbox/host"
	"os"
)

//go:generate go-bindata -dev -o assets.go -pkg gearbox admin/dist/...

func main() {
	gearbox.Instance = gearbox.NewApp(&gearbox.Args{
		HostConnector: host.GetConnector(),
		GlobalOptions: cmd.GlobalOptions,
	})
	status := gearbox.Instance.Initialize()
	if status.IsError() {
		fmt.Println(status.Message)
		if status.CliHelp != "" {
			fmt.Println(status.CliHelp)
		}
		os.Exit(1)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
