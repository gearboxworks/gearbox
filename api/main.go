package main

import (
	"gearbox"
	"gearbox/host"
)

func main() {
	c := gearbox.NewConfig(&host.MacOsConnector{})
	c.Initialize()
	gearbox.NewApi(c).Run()
}
