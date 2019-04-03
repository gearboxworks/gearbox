package routes

import "gearbox/gearbox"

type Gearbox struct {
	Gearbox gearbox.Gearboxer
	ProjectConnector
}

func NewGearbox(gb gearbox.Gearboxer) *Gearbox {
	return &Gearbox{
		Gearbox:          gb,
		ProjectConnector: ProjectConnector{Gearbox: gb},
	}
}
