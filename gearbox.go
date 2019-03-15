package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/host"
	"gearbox/only"
	"gearbox/stat"
	"log"
	"path/filepath"
)

const RepoRawBaseUrl = "https://raw.githubusercontent.com/gearboxworks/gearbox"

const DefaultAuthority = "gearbox.works"
const DefaultOrgName = "gearboxworks"

type JsonFileScope string

const (
	ProjectScope   JsonFileScope = "project"
	ContainerScope JsonFileScope = "container"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	StackMap      StackMap
	RequestType   api.ResourceName
	GlobalOptions *GlobalOptions
	hostApi       *HostApi
	errorLog      *ErrorLog
}

type Args Gearbox

func (me *Gearbox) Initialize() (status stat.Status) {
	for range only.Once {
		status = me.Config.Initialize()
		if status.IsError() {
			break
		}
		me.StackMap, status = me.GetStackMap()
		if status.IsError() {
			break
		}
	}
	return status
}

func NewGearbox(args *Args) *Gearbox {
	gb := Gearbox{
		HostConnector: args.HostConnector,
		GlobalOptions: args.GlobalOptions,
		Config:        args.Config,
		errorLog:      &ErrorLog{},
	}
	if args.Config == nil {
		gb.Config = NewConfig(&gb)
	}
	if args.GlobalOptions == nil {
		gb.GlobalOptions = &GlobalOptions{}
	}
	gb.hostApi = NewHostApi(&gb)
	return &gb
}

func (me *Gearbox) GetStackMap() (sm StackMap, status stat.Status) {
	for range only.Once {
		gears := NewGears(me)
		status = gears.Refresh()
		if status.IsError() {
			break
		}
		//gears.RoleMap
	}
	return sm, status
}

func (me *Gearbox) GetApiSelfLink(name api.ResourceName, vars api.UriTemplateVars) string {
	t, status := me.hostApi.GetApiSelfLink(name)
	if status.IsError() {
		// @TODO consider handling this with Status
		panic(status.Message)
	}
	return api.ExpandUriTemplate(t, vars)
}

func (me *Gearbox) FindProjectWithDetails(hostname string) (p *Project, status stat.Status) {
	return me.Config.Projects.FindProjectWithDetails(me, hostname)
}

func (me *Gearbox) GetProject(hostname string) (p *Project, status stat.Status) {
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

func (me *Gearbox) AddBasedir(dir string, nickname ...string) (status stat.Status) {
	status = me.Config.Basedirs.AddBasedir(me, dir, nickname...)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) UpdateBasedir(nickname string, dir string) (status stat.Status) {
	status = me.Config.Basedirs.UpdateBasedir(me, nickname, dir)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *Gearbox) DeleteNamedBasedir(nickname string) (status stat.Status) {
	status = me.Config.Basedirs.DeleteNamedBasedir(me, nickname)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}
func (me *Gearbox) ValidateBasedirNickname(nn string, args *validateArgs) stat.Status {
	args.Gearbox = me
	return ValidateBasedirNickname(nn, args)
}

func (me *Gearbox) ValidateProjectHostname(hn string, args *validateArgs) stat.Status {
	args.Gearbox = me
	return ValidateProjectHostname(hn, args)
}

func (me *Gearbox) RequestAvailableContainers(query ...*dockerhub.ContainerQuery) (names dockerhub.ContainerNames, status stat.Status) {
	for range only.Once {
		var _query *dockerhub.ContainerQuery
		if len(query) == 0 {
			_query = &dockerhub.ContainerQuery{}
		} else {
			_query = query[0]
		}
		dh := dockerhub.DockerHub{}
		names, status = dh.RequestAvailableContainerNames(_query)
	}
	return names, status
}

func getFirstBasedir(basedirs []string) (basedir string) {
	if len(basedirs) == 0 {
		basedir = PrimaryBasedirNickname
	} else {
		basedir = basedirs[0]
	}
	return basedir
}

func (me *Gearbox) GetProjectDir(path string, basedirs ...string) (basedir string, status stat.Status) {
	for range only.Once {
		var bd string
		bd, status = me.Config.GetHostBasedir(getFirstBasedir(basedirs))
		if status.IsError() {
			break
		}
		basedir = filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path))
	}
	return basedir, status
}

func (me *Gearbox) GetProjectFilepath(path string, basedirs ...string) (pfp string, status stat.Status) {
	for range only.Once {
		var pd string
		pd, status = me.GetProjectDir(path, getFirstBasedir(basedirs))
		if status.IsError() {
			break
		}
		pfp = filepath.FromSlash(fmt.Sprintf("%s/%s", pd, ProjectFilename))
	}
	return pfp, status
}

func (me *Gearbox) AddNamedStackToProject(stackName StackName, hostname string) (status stat.Status) {
	for range only.Once {
		var p *Project
		p, status = me.FindProjectWithDetails(hostname)
		if status.IsError() {
			break
		}
		status = p.AddNamedStack(stackName)
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("stack '%s' added to project '%s'", stackName, hostname)
	}
	return status
}

func (me *Gearbox) FindNamedStack(stackName StackName) (stack *Stack, status stat.Status) {
	stack, status = FindNamedStack(me, stackName)
	return stack, status
}
