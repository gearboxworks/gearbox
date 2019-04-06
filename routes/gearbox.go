package routes

import "gearbox/gearbox"

type Gearbox struct {
	Gearbox gearbox.Gearboxer
	ProjectModel
}

func NewGearbox(gb gearbox.Gearboxer) *Gearbox {
	return &Gearbox{
		Gearbox:      gb,
		ProjectModel: ProjectModel{Gearbox: gb},
	}
}
