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
	gearbox.Instance = gearbox.NewGearbox(&gearbox.GearboxArgs{
		HostConnector: &host.MacOsConnector{},
	})
	gearbox.Instance.Initialize()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
