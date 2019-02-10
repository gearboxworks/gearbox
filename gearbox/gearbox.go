package gearbox

import (
	"gearbox/gearbox/host"
	"log"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
}

func (me *Gearbox) Initialize() {
	hc := me.HostConnector
	if hc == nil {
		log.Fatal("Gearbox has no host connector. (End users should never see this; it is a programming error.)")
	}
	WriteAdminAssetsToWebRoot(hc)
	me.Config.Initialize()
}

func NewGearbox(hc host.Connector) *Gearbox {
	gb := Gearbox{
		HostConnector: hc,
		Config:        NewConfig(hc),
	}
	return &gb
}

func (me *Gearbox) AddProjectRoot(dir string) {
	pr := NewProjectRoot(me.Config.VmProjectRoot, dir)
	me.Config.ProjectRoots = append(me.Config.ProjectRoots, pr)
	me.Config.LoadProjectsAndWrite()
}
