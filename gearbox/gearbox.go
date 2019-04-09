package gearbox

import (
	"encoding/json"
	"gearbox/api"
	"gearbox/box"
	"gearbox/config"
	"gearbox/dockerhub"
	"gearbox/gears"
	"gearbox/global"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/project"
	"gearbox/ssh"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

type JsonFileScope string

const (
	ProjectScope   JsonFileScope = "project"
	ContainerScope JsonFileScope = "container"
)

var _ Gearboxer = (*Gearbox)(nil)

var Instance Gearboxer

type Gearboxer interface {
	//AddNamedStackToProject(gears.StackId, types.Hostname) status.Status
	//GetNamedStackNames() (Stacknames, status.Status)
	//GetProjectFilepath(string, string) (string, status.Status)
	//GetProjectWithDetails(config.Hostname) (*config.Project, status.Status)
	//ValidateBasedirNickname(string, *config.ValidateArgs) status.Status
	AddBasedir(types.AbsoluteDir, ...types.Nickname) status.Status
	AddNamedStack(*gears.NamedStack) status.Status
	AddProject(*project.Project) status.Status
	Admin(ViewerType)
	ConnectSSH(ssh.Args) status.Status
	CreateBox(box.Args) status.Status
	DeleteNamedBasedir(types.Nickname) status.Status
	DeleteNamedStack(stackid types.StackId) status.Status
	DeleteProject(hostname types.Hostname) status.Status
	FindNamedStack(stackid types.StackId) (*gears.NamedStack, status.Status)
	FindProject(hostname types.Hostname) (*project.Project, status.Status)
	GetConfig() config.Configer
	GetGears() *gears.Gears
	GetGlobalOptions() *global.Options
	GetApi() api.Apier
	GetNamedStackMap() (gears.NamedStackMap, status.Status)
	GetNamedStackRoleMap(types.StackId) (gears.StackRoleMap, status.Status)
	GetOsSupport() oss.OsSupporter
	GetProjectMap() (project.Map, status.Status)
	GetRouteName() types.RouteName
	GetStackRoleMap() (gears.StackRoleMap, status.Status)
	Initialize() status.Status
	IsDebug() bool
	NamedBasedirExists(types.Nickname) bool
	NoCache() bool
	PrintBoxStatus(box.Args) status.Status
	ProjectExists(types.Hostname) (bool, status.Status)
	RequestAvailableContainers(...*dockerhub.ContainerQuery) (dockerhub.ContainerNames, status.Status)
	RestartBox(box.Args) status.Status
	SetConfig(config.Configer)
	SetApi(api api.Apier)
	SetRouteName(types.RouteName)
	StartBox(box.Args) status.Status
	StopBox(box.Args) status.Status
	UpdateBasedir(types.Nickname, types.AbsoluteDir) status.Status
	UpdateNamedStack(*gears.NamedStack) status.Status
	UpdateProject(*project.Project) status.Status
	WriteLog([]byte) (int, error)
}

type Gearbox struct {
	Config        config.Configer
	OsSupport     oss.OsSupporter
	StackMap      gears.NamedStackMap
	GlobalOptions *global.Options
	Api           api.Apier
	RouteName     types.RouteName
	Gears         *gears.Gears
	errorLog      *ErrorLog
}

type Args Gearbox

func NewGearbox(args *Args) Gearboxer {
	gb := Gearbox{
		OsSupport:     args.OsSupport,
		GlobalOptions: args.GlobalOptions,
		Config:        args.Config,
		errorLog:      &ErrorLog{},
	}
	if args.Config == nil {
		gb.Config = config.NewConfig(args.OsSupport)
	}
	if args.GlobalOptions == nil {
		gb.GlobalOptions = &global.Options{}
	}
	if args.Api != nil {
		gb.Api = args.Api
		gb.Api.SetParent(&gb)
	}
	if args.Gears == nil {
		gb.Gears = gears.NewGears(gb.OsSupport)
	}
	return &gb
}

func (me *Gearbox) GetNamedStackMap() (nsm gears.NamedStackMap, sts status.Status) {
	for range only.Once {
		// @TODO Remove these comments after debugging
		//if me.StackMap != nil {
		//	break
		//}
		me.StackMap, sts = me.Gears.GetNamedStackMap()
		if is.Error(sts) {
			break
		}
	}
	return me.StackMap, sts
}

func (me *Gearbox) AddNamedStack(*gears.NamedStack) status.Status {
	panic("implement me")
}

func (me *Gearbox) UpdateNamedStack(*gears.NamedStack) status.Status {
	panic("implement me")
}

func (me *Gearbox) DeleteNamedStack(stackid types.StackId) status.Status {
	panic("implement me")
}

func (me *Gearbox) FindNamedStack(stackid types.StackId) (stack *gears.NamedStack, sts status.Status) {
	return me.Gears.FindNamedStack(stackid)
}

func (me *Gearbox) FindProject(hostname types.Hostname) (pp *project.Project, sts status.Status) {
	for range only.Once {
		var cp *config.Project
		cp, sts = me.Config.FindProject(hostname)
		if is.Error(sts) {
			break
		}
		pp = project.NewProject(cp)
		sts = pp.Renew(cp.Path)
		if is.Error(sts) {
			break
		}
	}
	return pp, sts
}

func (me *Gearbox) AddProject(p *project.Project) (sts status.Status) {
	// @TODO Also need to add project file
	return me.Config.AddProject(p.Project)
}
func (me *Gearbox) UpdateProject(p *project.Project) (sts status.Status) {
	// @TODO Also need to update project file
	return me.Config.UpdateProject(p.Project)
}
func (me *Gearbox) DeleteProject(hostname types.Hostname) (sts status.Status) {
	return me.Config.DeleteProject(hostname)
}
func (me *Gearbox) GetProjectMap() (pm project.Map, sts status.Status) {
	for range only.Once {
		var cpm config.ProjectMap
		cpm, sts = me.Config.GetProjectMap()
		if is.Error(sts) {
			break
		}
		pm = make(project.Map, 0)
		for i, cp := range cpm {
			pp := project.NewProject(cp)
			sts = pp.Renew(cp.Path)
			if is.Error(sts) {
				break
			}
			pm[i] = pp
		}
	}
	return pm, sts
}

func (me *Gearbox) GetStackRoleMap() (gears.StackRoleMap, status.Status) {
	return me.Gears.GetStackRoleMap()
}
func (me *Gearbox) GetNamedStackRoleMap(stackid types.StackId) (gears.StackRoleMap, status.Status) {
	return me.Gears.GetNamedStackRoleMap(stackid)
}

func (me *Gearbox) GetNamedStackIds() (types.StackIds, status.Status) {
	return me.Gears.GetNamedStackIds()
}

func (me *Gearbox) WriteLog(msg []byte) (int, error) {
	return me.errorLog.Write(msg)
}
func (me *Gearbox) GetGlobalOptions() *global.Options {
	return me.GlobalOptions
}

func (me *Gearbox) GetGears() *gears.Gears {
	return me.Gears
}

func (me *Gearbox) GetRouteName() types.RouteName {
	return me.RouteName
}

func (me *Gearbox) SetApi(api api.Apier) {
	me.Api = api
}

func (me *Gearbox) SetRouteName(routeName types.RouteName) {
	me.RouteName = routeName
}

func (me *Gearbox) GetOsSupport() oss.OsSupporter {
	return me.OsSupport
}

func (me *Gearbox) GetApi() api.Apier {
	return me.Api
}

func (me *Gearbox) GetConfig() config.Configer {
	return me.Config
}
func (me *Gearbox) SetConfig(config config.Configer) {
	me.Config = config
}

func (me *Gearbox) Initialize() (sts status.Status) {
	for range only.Once {
		sts = me.Gears.Initialize()
		if is.Error(sts) {
			break
		}
		sts = me.Config.Initialize()
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me *Gearbox) GetProject(hostname types.Hostname) (p *project.Project, sts status.Status) {
	for range only.Once {
		pm, sts := me.GetProjectMap()
		if status.IsError(sts) {
			break
		}
		p, sts = pm.GetProject(hostname)
	}
	return p, sts
}

//func (me Gearboxer) GetProjects() string {
//	j, err := json.Marshal(me.GetProjectMap())
//	if err != nil {
//		log.Fatal(err)
//	}
//	return string(j)
//}
//
func (me *Gearbox) Admin(viewer ViewerType) {
	aui := NewAdminUi(me, viewer)
	aui.Initialize()
	defer aui.Close()
	aui.Start()
}

func (me *Gearbox) ProjectExists(hostname types.Hostname) (ok bool, sts status.Status) {
	for range only.Once {
		pm, sts := me.GetProjectMap()
		if status.IsError(sts) {
			break
		}
		ok = pm.ProjectExists(hostname)
	}
	return ok, sts
}

func (me *Gearbox) NamedBasedirExists(nickname types.Nickname) bool {
	return me.Config.GetBasedirMap().NamedBasedirExists(nickname)
}

func (me *Gearbox) BasedirExists(dir types.AbsoluteDir) bool {
	return me.Config.GetBasedirMap().BasedirExists(dir)
}

func (me *Gearbox) AddBasedir(dir types.AbsoluteDir, nickname ...types.Nickname) (sts status.Status) {
	for range only.Once {
		sts = me.Config.AddBasedir(dir, nickname...)
		if is.Error(sts) {
			break
		}
		sts := me.Config.LoadProjectsAndWrite()
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me *Gearbox) UpdateBasedir(nickname types.Nickname, dir types.AbsoluteDir) (sts status.Status) {
	for range only.Once {
		sts = me.Config.GetBasedirMap().UpdateBasedir(nickname, dir)
		if status.IsError(sts) {
			break
		}
		sts = me.Config.LoadProjectsAndWrite()
		if status.IsError(sts) {
			break
		}
	}
	return sts
}

func (me *Gearbox) DeleteNamedBasedir(nickname types.Nickname) (sts status.Status) {
	for range only.Once {
		sts = me.Config.GetBasedirMap().DeleteNamedBasedir(nickname)
		if status.IsError(sts) {
			break
		}
		sts := me.Config.LoadProjectsAndWrite()
		if status.IsError(sts) {
			break
		}
	}
	return sts
}

func (me *Gearbox) RequestAvailableContainers(query ...*dockerhub.ContainerQuery) (names dockerhub.ContainerNames, sts status.Status) {
	for range only.Once {
		var _query *dockerhub.ContainerQuery
		if len(query) == 0 {
			_query = &dockerhub.ContainerQuery{}
		} else {
			_query = query[0]
		}
		dh := dockerhub.DockerHub{}
		names, sts = dh.RequestAvailableContainerNames(_query)
	}
	return names, sts
}

//func (me *Parent) GetProjectDir(path types.RelativePath, basedir types.Nickname) (bd types.AbsoluteDir, sts status.Status) {
//	for range only.Once {
//		var bd types.AbsoluteDir
//		bd, sts = me.Config.GetHostBasedir(basedir)
//		if status.IsError(sts) {
//			break
//		}
//		bd = types.AbsoluteDir(filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path)))
//	}
//	return bd, sts
//}

//func (me *Parent) GetProjectFilepath(path types.RelativePath, basedir types.Nickname) (pfp types.AbsoluteDir, sts status.Status) {
//	for range only.Once {
//		var pd types.AbsoluteDir
//		pd, sts = me.GetProjectDir(path, basedir)
//		if status.IsError(sts) {
//			break
//		}
//		pfp = types.AbsoluteDir(filepath.FromSlash(fmt.Sprintf("%s/%s", pd, jsonfile.BaseFilename)))
//	}
//	return pfp, sts
//}

//func (me *Parent) AddNamedStackToProject(stackid gears.StackId, hostname types.Hostname) (sts status.Status) {
//	for range only.Once {
//		var p *config.Project
//		p, sts = me.GetProjects(hostname)
//		if status.IsError(sts) {
//			break
//		}
//		sts = p.AddNamedStack(stackid)
//		if status.IsError(sts) {
//			break
//		}
//		sts = status.Success("named stack ID '%s' added to project '%s'", stackid, hostname)
//	}
//	return sts
//}

func (me *Gearbox) ConnectSSH(sshArgs ssh.Args) (sts status.Status) {
	return ssh.NewSSH(sshArgs).StartSSH()
}

//
// This just here as a method to copy when needed
//
func (me *Gearbox) Clone() *Gearbox {
	gb := Gearbox{}
	for range only.Once {
		b, err := json.Marshal(me)
		if err != nil {
			break
		}
		_ = json.Unmarshal(b, &gb)
	}
	return &gb
}
