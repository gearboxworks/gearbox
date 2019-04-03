package main

import (
	"fmt"
	"gearbox/app/cmd"
	"gearbox/gearbox"
	"gearbox/hostapi"
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
	ha := hostapi.NewHostApi(gb)
	sts := ha.Route()
	if is.Error(sts) {
		panic(sts.Message())
	}
	gb.SetHostApi(ha)
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
