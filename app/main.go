package main

import (
	"fmt"
	"gearbox/api"
	"gearbox/apimvc"
	"gearbox/app/cmd"
	"gearbox/app/logger"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"log"
	"net/http"
	"os"
)

//go:generate go-bindata -dev -o gearbox/assets.go -pkg gearbox admin/dist/...

var startProfiler = false

func main() {

	for range only.Once {

		if startProfiler == true {
			go func() {
				log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
			}()
		}

		var sts status.Status

		_, sts = logger.NewLogger(oss.Get(), logger.Logger{})
		if is.Error(sts) {
			break
		}

		sts = status.Success("gearbox started")
		status.Log(sts)

		gb := gearbox.NewGearbox(&gearbox.Args{
			OsSupport:     oss.Get(),
			GlobalOptions: cmd.GlobalOptions,
		})

		gearbox.Instance = gb
		gb.SetApi(api.NewApi(gb))
		sts = apimvc.AddControllers(gb)
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
}

