package apimvc

import (
	"gearbox/gearbox"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

func noop(i ...interface{}) interface{} {
	return nil
}

func AddControllers(gb gearbox.Gearboxer) (sts Status) {
	for range only.Once {

		controllers := []ListController{
			NewDirectoryController(),
			NewProjectController(gb),
			NewStackController(gb),
			NewServiceController(gb),
			NewGearspecController(gb),
			NewAuthorityController(gb),
			NewBasedirController(gb.GetConfig()),
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
