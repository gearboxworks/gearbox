package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/box"
	"gearbox/config"
	"gearbox/dockerhub"
	"gearbox/gears"
	"gearbox/global"
	"gearbox/project"
	"gearbox/service"
	"gearbox/ssh"
	"gearbox/types"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"log"
	"path/filepath"
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
	AddBasedir(types.Dir, ...types.Nickname) status.Status
	AddNamedStack(*gears.NamedStack) status.Status
	AddProject(*project.Project) status.Status
	Admin(ViewerType)
	DeleteBasedir(types.Nickname) status.Status
	DeleteNamedStack(stackid types.StackId) status.Status
	DeleteProject(hostname types.Hostname) status.Status
	FindNamedStack(stackid types.StackId) (*gears.NamedStack, status.Status)
	FindService(serviceid service.Identifier) (*gears.Service, status.Status)
	FindProject(hostname types.Hostname) (*project.Project, status.Status)
	GetConfig() config.Configer
	GetGears() *gears.Gears
	GetGlobalOptions() *global.Options
	GetApi() api.Apier
	GetNamedStackMap() (gears.NamedStackMap, status.Status)
	GetServiceMap() (gears.ServiceMap, status.Status)
	GetNamedStackRoleMap(types.StackId) (gears.StackRoleMap, status.Status)
	GetOsBridge() osbridge.OsBridger
	GetProjectMap() (project.Map, status.Status)
	GetRouteName() types.RouteName
	GetStackRoleMap() (gears.StackRoleMap, status.Status)
	Initialize() status.Status
	IsDebug() bool
	BasedirExists(types.Nickname) bool
	NoCache() bool
	ProjectExists(types.Hostname) (bool, status.Status)
	RequestAvailableContainers(...*dockerhub.ContainerQuery) (dockerhub.ContainerNames, status.Status)

	// VM related.
	RunAsDaemon(*box.Args) status.Status
	StartBox(*box.Args) status.Status
	StopBox(*box.Args) status.Status
	RestartBox(*box.Args) status.Status
	CreateBox(*box.Args) status.Status
	PrintBoxStatus(*box.Args) status.Status
	ConnectSSH(ssh.Args) status.Status

	SetConfig(config.Configer)
	SetApi(api api.Apier)
	SetRouteName(types.RouteName)
	UpdateBasedir(types.Nickname, types.Dir) status.Status
	UpdateNamedStack(*gears.NamedStack) status.Status
	UpdateProject(*project.Project) status.Status
	WriteLog([]byte) (int, error)
}

type Gearbox struct {
	Config        config.Configer
	OsBridge      osbridge.OsBridger
	StackMap      gears.NamedStackMap
	ServiceMap    gears.ServiceMap
	GlobalOptions *global.Options
	Api           api.Apier
	RouteName     types.RouteName
	Gears         *gears.Gears
	errorLog      *ErrorLog
}

type Args Gearbox

func NewGearbox(args *Args) Gearboxer {
	gb := Gearbox{
		OsBridge:      args.OsBridge,
		GlobalOptions: args.GlobalOptions,
		Config:        args.Config,
		errorLog:      &ErrorLog{},
	}
	if args.Config == nil {
		gb.Config = config.NewConfig(args.OsBridge)
	}
	if args.GlobalOptions == nil {
		gb.GlobalOptions = &global.Options{}
	}
	if args.Api != nil {
		gb.Api = args.Api
		gb.Api.SetParent(&gb)
	}
	if args.Gears == nil {
		gb.Gears = gears.NewGears(gb.OsBridge)
	}
	return &gb
}

func (me *Gearbox) GetServiceMap() (nsm gears.ServiceMap, sts status.Status) {
	for range only.Once {
		// @TODO Remove these comments after debugging
		//if me.ServiceMap != nil {
		//	break
		//}
		me.ServiceMap, sts = me.Gears.GetServiceMap()
		if is.Error(sts) {
			break
		}
	}
	return me.ServiceMap, sts
}

func (me *Gearbox) GetNamedStackMap() (nsm gears.NamedStackMap, sts status.Status) {
	for range only.Once {
		if me.StackMap != nil {
			break
		}
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

func (me *Gearbox) FindService(serviceid service.Identifier) (service *gears.Service, sts status.Status) {
	return me.Gears.FindService(serviceid)
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

func (me *Gearbox) GetOsBridge() osbridge.OsBridger {
	return me.OsBridge
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
		me.WriteAssetsToAdminWebRoot()

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

func (me *Gearbox) WriteAssetsToAdminWebRoot() {
	hc := me.OsBridge
	if hc == nil {
		log.Fatal("Gearbox has no osbridge connector. (End users should never see this; it is a programming error.)")
	}
	for _, afn := range AssetNames() {
		afn = filepath.FromSlash(afn)
		err := RestoreAsset(hc.GetUserConfigDir(), afn)
		if err != nil {
			afn = fmt.Sprintf("'%s/%s'", hc.GetUserConfigDir(), afn)
			log.Printf("Could not restore asset '%s': %v\n", afn, err.Error())
		}
	}

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

func (me *Gearbox) BasedirExists(nickname types.Nickname) bool {
	return me.Config.GetBasedirMap().NicknameExists(nickname)
}

func (me *Gearbox) AddBasedir(basedir types.Dir, nickname ...types.Nickname) (sts status.Status) {
	var nn types.Nickname
	if len(nickname) == 0 {
		nn = ""
	}
	for range only.Once {
		_, sts = me.Config.AddBasedir(&config.BasedirArgs{
			Basedir:  basedir,
			Nickname: nn,
		})
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

func (me *Gearbox) UpdateBasedir(nickname types.Nickname, dir types.Dir) (sts status.Status) {
	for range only.Once {
		sts = me.Config.GetBasedirMap().UpdateBasedir(
			me.Config,
			config.NewBasedir(nickname, dir),
		)
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

func (me *Gearbox) DeleteBasedir(nickname types.Nickname) (sts status.Status) {
	for range only.Once {
		sts = me.Config.DeleteBasedir(nickname)
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

//func (me *Parent) GetProjectDir(path types.Path, basedir types.Nickname) (bd types.Dir, sts status.Status) {
//	for range only.Once {
//		var bd types.Dir
//		bd, sts = me.Config.GetBasedir(basedir)
//		if status.IsError(sts) {
//			break
//		}
//		bd = types.Dir(filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path)))
//	}
//	return bd, sts
//}

//func (me *Parent) GetProjectFilepath(path types.Path, basedir types.Nickname) (pfp types.Dir, sts status.Status) {
//	for range only.Once {
//		var pd types.Dir
//		pd, sts = me.GetProjectDir(path, basedir)
//		if status.IsError(sts) {
//			break
//		}
//		pfp = types.Dir(filepath.FromSlash(fmt.Sprintf("%s/%s", pd, jsonfile.BaseFilename)))
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
func (me *Gearbox) IsDebug() bool {
	return me.GlobalOptions.IsDebug
}

func (me *Gearbox) NoCache() bool {
	return me.GlobalOptions.NoCache
}
