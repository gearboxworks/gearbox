package apimvc

import (
	"gearbox/api"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

func AddControllers(a api.Apier, gb gearbox.Gearboxer) (sts status.Status) {
	for range only.Once {

		controllers := []apimodeler.ListController{
			NewProjectController(gb),
			NewStackController(gb),
			NewServiceController(gb),
			NewGearspecController(gb),
			NewAuthorityController(gb),
			NewRootController(gb),
		}

		for _, cs := range controllers {
			sts = a.AddController(cs)
			if is.Error(sts) {
				panic(sts.Message())
			}
		}

	}
	return sts
}
