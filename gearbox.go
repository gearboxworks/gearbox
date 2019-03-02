package gearbox

import (
	"encoding/json"
	"gearbox/dockerhub"
	"gearbox/host"
	"log"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	Stacks        StackMap
}

type GearboxArgs Gearbox

func (me *Gearbox) Initialize() {
	me.Config.Initialize()
}

func NewGearbox(args *GearboxArgs) *Gearbox {
	if args.Config == nil {
		args.Config = NewConfig(args.HostConnector)
	}
	gb := Gearbox{
		HostConnector: args.HostConnector,
		Config:        args.Config,
		Stacks:        GetStackMap(),
	}
	return &gb
}

func (me *Gearbox) GetProjects() string {
	j, err := json.Marshal(me.Config.Projects)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func (me *Gearbox) Admin(viewer ViewerType) {
	aui := NewAdminUi(me, viewer)
	aui.Initialize()
	defer aui.Close()
	aui.Start()
}

func (me *Gearbox) AddProjectRoot(dir string) {
	pr := NewProjectRoot(me.Config.VmProjectRoot, dir)
	me.Config.ProjectRoots = append(me.Config.ProjectRoots, pr)
	me.Config.LoadProjectsAndWrite()
}

func (me *Gearbox) StartVm() {
	vm := &Vm{}
	err := vm.StartVm()
	if err != nil {
		panic(err)
	}
	return
}

func (me *Gearbox) RequestAvailableContainers(query ...*dockerhub.ContainerQuery) dockerhub.ContainerNames {
	var _query *dockerhub.ContainerQuery
	if len(query) == 0 {
		_query = &dockerhub.ContainerQuery{}
	} else {
		_query = query[0]
	}
	dh := dockerhub.DockerHub{}
	return dh.RequestAvailableContainerNames(_query)
}
