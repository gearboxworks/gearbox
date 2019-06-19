package main

import (
	"fmt"
	"gearbox/api"
	"gearbox/apimvc"
	"gearbox/app/cmd"
	"gearbox/gearbox"
	"gearbox/global"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"os"
)

//go:generate go-bindata -dev -o gearbox/dist.go -pkg gearbox app/dist/... admin/dist/...

func main() {
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsBridge:      gearbox.GetOsBridge(global.Brandname, global.UserDataPath),
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
