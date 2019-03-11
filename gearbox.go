package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/dockerhub"
	"gearbox/host"
	"gearbox/only"
	"log"
	"path/filepath"
)

type JsonFileScope string

const (
	ProjectScope   JsonFileScope = "project"
	ContainerScope JsonFileScope = "container"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	Stacks        StackMap
	RequestType   string
	Options       *GlobalOptions
	hostApi       *HostApi
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
	gb.hostApi = NewHostApi(&gb)
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

func (me *Gearbox) NamedBasedirExists(nickname string) bool {
	return me.Config.Basedirs.NamedBasedirExists(nickname)
}

func (me *Gearbox) BasedirExists(dir string) bool {
	return me.Config.Basedirs.BasedirExists(dir)
}

func (me *Gearbox) AddBasedir(dir string, nickname ...string) (status Status) {
	status = me.Config.Basedirs.AddBasedir(me, dir, nickname...)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) UpdateBasedir(nickname string, dir string) (status Status) {
	status = me.Config.Basedirs.UpdateBasedir(me, nickname, dir)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) DeleteNamedBasedir(nickname string) (status Status) {
	status = me.Config.Basedirs.DeleteNamedBasedir(me, nickname)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}
func (me *Gearbox) ValidateBasedirNickname(nn string, args *validateArgs) Status {
	args.Gearbox = me
	return ValidateBasedirNickname(nn, args)
}

func (me *Gearbox) ValidateProjectHostname(hn string, args *validateArgs) Status {
	args.Gearbox = me
	return ValidateProjectHostname(hn, args)
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

func getFirstBasedir(basedirs []string) (basedir string) {
	if len(basedirs) == 0 {
		basedir = PrimaryBasedirNickname
	} else {
		basedir = basedirs[0]
	}
	return basedir
}

func (me *Gearbox) GetProjectDir(path string, basedirs ...string) (basedir string, err error) {
	for range only.Once {
		var bd string
		bd, err = me.Config.GetHostBasedir(getFirstBasedir(basedirs))
		if err != nil {
			break
		}
		basedir = filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path))
	}
	return basedir, err
}

func (me *Gearbox) GetProjectFilepath(path string, basedirs ...string) (pfp string, err error) {
	for range only.Once {
		var pd string
		pd, err = me.GetProjectDir(path, getFirstBasedir(basedirs))
		if err != nil {
			break
		}
		pfp = filepath.FromSlash(fmt.Sprintf("%s/%s", pd, ProjectFilename))
	}
	return pfp, err
}
