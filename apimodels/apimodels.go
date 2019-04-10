package apimodels

import (
	"gearbox/api"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

func AddModels(a api.Apier, gb gearbox.Gearboxer) (sts status.Status) {
	for range only.Once {

		sts = a.AddModels(NewProjectModel(gb))
		if is.Error(sts) {
			panic(sts.Message())
		}

		sts = a.AddModels(NewStackModel(gb))
		if is.Error(sts) {
			panic(sts.Message())
		}

		sts = a.AddModels(NewRootModel(gb))
		if is.Error(sts) {
			panic(sts.Message())
		}

	}
	return sts
}
