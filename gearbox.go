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
	RequestType   string
	Options       *GlobalOptions
}

type GearboxArgs Gearbox

func (me *Gearbox) Initialize() (status Status) {
	return me.Config.Initialize()
}

func NewGearbox(args *GearboxArgs) *Gearbox {
	gb := Gearbox{
		HostConnector: args.HostConnector,
		Options:       args.Options,
		Config:        args.Config,
		Stacks:        GetStackMap(),
	}
	if args.Config == nil {
		gb.Config = NewConfig(&gb)
	}
	if args.Options == nil {
		gb.Options = &GlobalOptions{}
	}
	return &gb
}

func (me *Gearbox) GetProjectResponse(hostname string) (pr *ProjectResponse, status Status) {
	return me.Config.Projects.GetProjectResponse(me, hostname)
}

func (me *Gearbox) GetProject(hostname string) (p *Project, status Status) {
	return me.Config.Projects.GetProject(me, hostname)
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

func (me *Gearbox) ProjectExists(hostname string) (ok bool) {
	return me.Config.Projects.ProjectExists(hostname)
}

func (me *Gearbox) NamedBaseDirExists(nickname string) bool {
	return me.Config.BaseDirs.NamedBaseDirExists(nickname)
}

func (me *Gearbox) BaseDirExists(dir string) bool {
	return me.Config.BaseDirs.BaseDirExists(dir)
}

func (me *Gearbox) AddBaseDir(dir string, nickname ...string) (status Status) {
	status = me.Config.BaseDirs.AddBaseDir(me, dir, nickname...)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) UpdateBaseDir(nickname string, dir string) (status Status) {
	status = me.Config.BaseDirs.UpdateBaseDir(me, nickname, dir)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) DeleteNamedBaseDir(nickname string) (status Status) {
	status = me.Config.BaseDirs.DeleteNamedBaseDir(me, nickname)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}
func (me *Gearbox) ValidateBaseDirNickname(nn string, args *validateArgs) Status {
	args.Gearbox = me
	return ValidateBaseDirNickname(nn, args)
}

func (me *Gearbox) ValidateProjectHostname(hn string, args *validateArgs) Status {
	args.Gearbox = me
	return ValidateProjectHostname(hn, args)
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
