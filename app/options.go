package main

import (
	"fmt"
	"gearbox"
	"gearbox/os_support"
	"log"
)

func main() {
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsSupport: oss.Get(),
	})

	gears := gb.GetGears()
	sts := gears.Initialize()
	if status.IsError(sts) {
		log.Fatal(status.Message)
	}

	fmt.Println(gears.String())

}
