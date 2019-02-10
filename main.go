package main

import (
	"fmt"
	"gearbox/cmd"
	"gearbox/gearbox"
	"gearbox/gearbox/host"
	"os"
)

//go:generate go-bindata -dev -o gearbox/asset.go -pkg gearbox admin/...

func main() {
	gearbox.Instance = gearbox.NewGearbox(&host.MacOsConnector{})
	gearbox.Instance.Initialize()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
