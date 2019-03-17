package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/host"
	"gearbox/only"
	"gearbox/stat"
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

var Instance Gearbox

type Gearbox interface {
	Admin(ViewerType)
	StartBox(BoxArgs) error
	StopBox(BoxArgs) error
	PrintBoxStatus(BoxArgs) (string, error)
	RestartBox(BoxArgs) error
	CreateBox(BoxArgs) (string, error)
	ConnectSSH(SSHArgs) error
	Initialize() stat.Status
	GetConfig() Config
	SetConfig(Config)
	GetHostConnector() host.Connector
	GetStackMap() (StackMap, stat.Status)
	GetGlobalOptions() *GlobalOptions
	GetHostApi() *HostApi
	GetResourceName() api.ResourceName
	SetResourceName(api.ResourceName)
	IsDebug() bool
	NoCache() bool
	ProjectExists(string) bool
	ValidateBasedirNickname(string, *ValidateArgs) stat.Status
	AddBasedir(string, ...string) stat.Status
	UpdateBasedir(string, string) stat.Status
	DeleteNamedBasedir(string) stat.Status
	NamedBasedirExists(string) bool
	FindProjectWithDetails(string) (*Project, stat.Status)
	AddNamedStackToProject(StackName, string) stat.Status
	RequestAvailableContainers(...*dockerhub.ContainerQuery) (dockerhub.ContainerNames, stat.Status)
	GetApiUrl(api.ResourceName, api.UriTemplateVars) (string, stat.Status)
	GetProjectFilepath(string, string) (string, stat.Status)
	WriteLog([]byte) (int, error)
}

type GearboxObj struct {
	Config        Config
	HostConnector host.Connector
	StackMap      StackMap
	GlobalOptions *GlobalOptions
	HostApi       *HostApi
	ResourceType  api.ResourceName
	errorLog      *ErrorLog
}

type Args GearboxObj

func (me *GearboxObj) WriteLog(msg []byte) (nn int, err error) {
	return me.errorLog.Write(msg)
}
func (me *GearboxObj) GetGlobalOptions() *GlobalOptions {
	return me.GlobalOptions
}

func (me *GearboxObj) GetResourceName() api.ResourceName {
	return me.ResourceType
}

func (me *GearboxObj) SetResourceName(resourceName api.ResourceName) {
	me.ResourceType = resourceName
}

func (me *GearboxObj) GetHostConnector() host.Connector {
	return me.HostConnector
}

func (me *GearboxObj) GetHostApi() *HostApi {
	return me.HostApi
}

func (me *GearboxObj) GetConfig() Config {
	return me.Config
}
func (me *GearboxObj) SetConfig(config Config) {
	me.Config = config
}

func (me *GearboxObj) Initialize() (status stat.Status) {
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

func NewApp(args *Args) Gearbox {
	gb := GearboxObj{
		HostConnector: args.HostConnector,
		GlobalOptions: args.GlobalOptions,
		Config:        args.Config,
		errorLog:      &ErrorLog{},
	}
	if args.Config == nil {
		gb.Config = NewConfiguration(&gb)
	}
	if args.GlobalOptions == nil {
		gb.GlobalOptions = &GlobalOptions{}
	}
	gb.HostApi = NewHostApi(&gb)
	return &gb
}

func (me *GearboxObj) GetStackMap() (sm StackMap, status stat.Status) {
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

func (me *GearboxObj) GetApiUrl(name api.ResourceName, vars api.UriTemplateVars) (url string, status stat.Status) {
	return me.HostApi.GetUrl(name, vars)
}

func (me *GearboxObj) FindProjectWithDetails(hostname string) (p *Project, status stat.Status) {
	return me.Config.GetProjectMap().FindProjectWithDetails(me, hostname)
}

func (me *GearboxObj) GetProject(hostname string) (p *Project, status stat.Status) {
	return me.Config.GetProjectMap().GetProject(me, hostname)
}

//func (me Gearbox) GetProjects() string {
//	j, err := json.Marshal(me.Config.GetProjectMap())
//	if err != nil {
//		log.Fatal(err)
//	}
//	return string(j)
//}
//
func (me *GearboxObj) Admin(viewer ViewerType) {
	aui := NewAdminUi(me, viewer)
	aui.Initialize()
	defer aui.Close()
	aui.Start()
}

func (me *GearboxObj) ProjectExists(hostname string) (ok bool) {
	return me.Config.GetProjectMap().ProjectExists(hostname)
}

func (me *GearboxObj) NamedBasedirExists(nickname string) bool {
	return me.Config.GetBasedirMap().NamedBasedirExists(nickname)
}

func (me *GearboxObj) BasedirExists(dir string) bool {
	return me.Config.GetBasedirMap().BasedirExists(dir)
}

func (me *GearboxObj) AddBasedir(dir string, nickname ...string) (status stat.Status) {
	status = me.Config.GetBasedirMap().AddBasedir(me, dir, nickname...)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *GearboxObj) UpdateBasedir(nickname string, dir string) (status stat.Status) {
	status = me.Config.GetBasedirMap().UpdateBasedir(me, nickname, dir)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}

func (me *GearboxObj) DeleteNamedBasedir(nickname string) (status stat.Status) {
	status = me.Config.GetBasedirMap().DeleteNamedBasedir(me, nickname)
	if !status.IsError() {
		status2 := me.Config.LoadProjectsAndWrite()
		if status2.IsError() {
			status = status2
		}
	}
	return status
}
func (me *GearboxObj) ValidateBasedirNickname(nn string, args *ValidateArgs) stat.Status {
	args.Gearbox = me
	return ValidateBasedirNickname(nn, args)
}

func (me *GearboxObj) ValidateProjectHostname(hn string, args *ValidateArgs) stat.Status {
	args.Gearbox = me
	return ValidateProjectHostname(hn, args)
}

func (me *GearboxObj) RequestAvailableContainers(query ...*dockerhub.ContainerQuery) (names dockerhub.ContainerNames, status stat.Status) {
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

func (me *GearboxObj) GetProjectDir(path string, basedir string) (bd string, status stat.Status) {
	for range only.Once {
		var bd string
		bd, status = me.Config.GetHostBasedir(basedir)
		if status.IsError() {
			break
		}
		basedir = filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path))
	}
	return basedir, status
}

func (me *GearboxObj) GetProjectFilepath(path string, basedir string) (pfp string, status stat.Status) {
	for range only.Once {
		var pd string
		pd, status = me.GetProjectDir(path, basedir)
		if status.IsError() {
			break
		}
		pfp = filepath.FromSlash(fmt.Sprintf("%s/%s", pd, ProjectFilename))
	}
	return pfp, status
}

func (me *GearboxObj) AddNamedStackToProject(stackName StackName, hostname string) (status stat.Status) {
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

func (me *GearboxObj) FindNamedStack(stackName StackName) (stack *Stack, status stat.Status) {
	stack, status = FindNamedStack(me, stackName)
	return stack, status
}
