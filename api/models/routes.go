package models

import (
	"gearbox/api"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

func AddRoutes(a api.Apier, gb gearbox.Gearboxer) (sts status.Status) {
	for range only.Once {
		sts = a.AddModels(NewProjectModel(gb))
		if is.Error(sts) {
			panic(sts.Message())
		}
		sts = a.AddModels(NewStackConnector(gb))
		if is.Error(sts) {
			panic(sts.Message())
		}
	}
	return sts
}
