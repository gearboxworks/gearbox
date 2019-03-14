package main

import (
	"fmt"
	"gearbox"
	"gearbox/host"
	"log"
)

func main() {
	gb := gearbox.NewGearbox(&gearbox.GearboxArgs{
		HostConnector: host.GetConnector(),
	})

	options := gearbox.NewOptions(gb)
	err := options.Refresh()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(options.String())

}
