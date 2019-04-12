package apimodels

import (
	"gearbox/api"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

func AddModels(a api.Apier, gb gearbox.Gearboxer) (sts status.Status) {
	for range only.Once {

		models := []apimodeler.ApiModeler{
			NewProjectModel(gb),
			NewStackModel(gb),
			NewServiceModel(gb),
			NewGearspecModel(gb),
			NewAuthorityModel(gb),
			NewRootModel(gb),
		}

		for _, ms := range models {
			sts = a.AddModels(ms)
			if is.Error(sts) {
				panic(sts.Message())
			}
		}

	}
	return sts
}
