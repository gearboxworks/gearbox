package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/dockerhub"
	"gearbox/host"
	"log"
	"net/http"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	Stacks        StackMap
}

type GearboxArgs Gearbox

func (me *Gearbox) Initialize() (status *Status) {
	return me.Config.Initialize()
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

func (me *Gearbox) AddBaseDir(dir string, nickname ...string) (status *Status) {
	var nn string
	if len(nickname) > 0 {
		nn = nickname[0]
	}
	bd := NewBaseDir(dir, &BaseDirArgs{
		VmDir:    me.Config.VmBaseDir,
		Nickname: nn,
	})
	if bd.Error != nil {
		status = NewStatus(&StatusArgs{
			HttpStatus: http.StatusBadRequest,
			Error:      bd.Error,
		})
		if dir == "" {
			status.Message = fmt.Sprint("value provide for base dir in 'host_dir' property was empty")
		} else {
			status.Message = fmt.Sprintf("could add add base dir '%s'; the ~ could not be expanded", dir)
		}
	}
	me.Config.BaseDirs[bd.Nickname] = bd
	status = me.Config.LoadProjectsAndWrite()
	return status
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
