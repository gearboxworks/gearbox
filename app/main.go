package main

import (
	"fmt"
	"gearbox/api"
	"gearbox/apimvc"
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
	gb.SetApi(api.NewApi(gb))
	sts := apimvc.AddControllers(gb)
	if is.Error(sts) {
		panic(sts.Message())
	}
	gb.GetApi().WireRoutes()
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
