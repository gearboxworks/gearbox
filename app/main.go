package main

import (
	"fmt"
	"gearbox/api"
	"gearbox/apimodels"
	"gearbox/app/cmd"
	"gearbox/gearbox"
	"gearbox/os_support"
	"gearbox/status"
	"gearbox/status/is"
	"os"
)

//go:generate go-bindata -dev -o assets.go -pkg gearbox admin/dist/...

func main() {
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsSupport:     oss.Get(),
		GlobalOptions: cmd.GlobalOptions,
	})
	gearbox.Instance = gb
	a := api.NewApi(gb)
	sts := apimodels.AddRoutes(a, gb)
	if is.Error(sts) {
		panic(sts.Message())
	}
	a.ConnectRoutes()
	gb.SetApi(a)
	sts = gb.Initialize()
	if status.IsError(sts) {
		fmt.Println(sts.Message())
		help := sts.GetHelp(status.CliHelp)
		if help != "" {
			fmt.Println(help)
		}
		os.Exit(1)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
