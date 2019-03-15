package main

import (
	"fmt"
	"gearbox"
	"gearbox/host"
	"log"
)

func main() {
	gb := gearbox.NewGearbox(&gearbox.Args{
		HostConnector: host.GetConnector(),
	})

	gears := gearbox.NewGears(gb)
	status := gears.Refresh()
	if status.IsError() {
		log.Fatal(status.Message)
	}

	fmt.Println(gears.String())

}
