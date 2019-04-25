package apimvc

import (
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status/is"
)

func noop(i ...interface{}) {}

func AddControllers(gb gearbox.Gearboxer) (sts Status) {
	for range only.Once {

		controllers := []ListController{
			NewProjectController(gb),
			NewStackController(gb),
			NewServiceController(gb),
			NewGearspecController(gb),
			NewAuthorityController(gb),
			NewRootController(gb),
		}
		a := gb.GetApi()
		for _, cs := range controllers {
			sts = a.AddController(cs)
			if is.Error(sts) {
				panic(sts.Message())
			}
		}

	}
	return sts
}
