package main

import (
	"gearbox"
	"gearbox/host"
)

func main() {
	c := gearbox.NewConfig(&host.WinConnector{})
	c.Initialize()
	gearbox.NewApi(c).Start()
}
