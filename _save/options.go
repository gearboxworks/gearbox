package _save

import (
	"fmt"
	"gearbox"
	//	"gearbox/os_support"
	"gearbox/status/is"
	"log"
)

func main_options() {
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsSupport: oss.Get(),
	})

	gears := gb.GetGears()
	sts := gears.Initialize()
	if is.Error(sts) {
		log.Fatal(sts.Message())
	}
	fmt.Println(gears.String())

}
