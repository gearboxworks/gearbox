package project

import (
	"fmt"
	"gearbox/config"
	"gearbox/gears"
	"gearbox/jsonfile"
	"gearbox/only"
	"gearbox/service"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/http"
	"path/filepath"
)

type Map map[types.Hostname]*Project
type Projects []*Project

type Project struct {
	loaded bool
	*config.Project
	Filepath types.Filepath      `json:"filepath"`
	Aliases  HostnameAliases     `json:"aliases"`
	Stack    service.ServicerMap `json:"stack"`
}

func NewProject(cp *config.Project) (p *Project) {
	p = &Project{
		Project: cp,

		//Stack: ConvertServiceMap(p.GetServicerMap()),
		//ServiceIds: ConvertServices(pp.GetServicerMap()),
		//
		//Aliases:    pp.GetAliases(),
	}
	return p
}

func (me Map) GetProject(hostname types.Hostname) (p *Project, sts status.Status) {
	for range only.Once {
		p, sts = me.GetProject(hostname)
		if status.IsError(sts) {
			break
		}
		sts = p.Load()
		if status.IsError(sts) {
			break
		}
		sts = status.Success("got project '%s'", hostname)
	}
	return p, sts
}

func (me Map) ProjectExists(hostname types.Hostname) (ok bool) {
	_, ok = me[hostname]
	return ok
}

func (me *Project) Renew(path types.Path) (sts status.Status) {
	for range only.Once {
		me.Path = path
		if me.Hostname == "" {
			me.Hostname = me.GetHostname()
		}
		if !me.loaded {
			sts = me.Load()
			if status.IsError(sts) {
				break
			}
		}
		for gs, ps := range me.Stack {
			gs, sts := gs.GetExpandedIdentifier()
			if is.Error(sts) {
				break
			}
			svc, sts := ps.GetServiceValue()
			if is.Error(sts) {
				sts = status.Fail().
					SetHttpStatus(http.StatusBadRequest).
					SetMessage("unable to get stack service value for gearspecid '%s'", gs).
					SetAllHelp("ensure your %s is using the correct roles.", gears.JsonFilename) // @TODO improve this help
				break
			}
			me.Stack[gs].Servicer = svc
		}
	}
	return sts
}

func (me *Project) Enable() (sts status.Status) {
	for range only.Once {

		// MICKMAKE - To be added.
		if is.Error(sts) {
			break
		}
	}

	return sts
}

func (me *Project) Disable() (sts status.Status) {
	for range only.Once {

		// MICKMAKE - To be added.
		if is.Error(sts) {
			break
		}
	}

	return sts
}

func (me *Project) Load() (sts status.Status) {
	for range only.Once {
		var fp types.Filepath
		fp, sts = me.GetFilepath()
		if is.Error(sts) {
			break
		}
		var j []byte
		j, sts = util.ReadBytes(fp)
		if is.Error(sts) {
			break
		}
		jpf := NewJsonFile(fp)
		if len(j) > 0 {
			sts = jpf.Unmarshal(j)
		}
		sts = jpf.FixupStack()
		if is.Error(sts) {
			break
		}
		me.Aliases = jpf.Aliases
		me.Stack = jpf.Stack
		me.Filepath = jpf.Filepath
		me.loaded = true
	}
	return sts
}

func (me *Project) GetAliases() (aliases HostnameAliases) {
	return me.Aliases
}

func (me *Project) GetServicerMap() (simap service.ServicerMap) {
	return me.Stack
}

func (me *Project) GetFilepath() (fp types.Filepath, sts status.Status) {
	var bd types.Dir
	for range only.Once {
		if me.Filepath != "" {
			break
		}
		bd, sts = me.Config.GetBasedir(me.Basedir)
		if is.Error(sts) {
			break
		}
		me.Filepath = types.Filepath(filepath.FromSlash(fmt.Sprintf("%s/%s/%s",
			bd,
			me.Path,
			jsonfile.BaseFilename,
		)))
	}
	return me.Filepath, sts
}

func (me *Project) GetProjectDir() (dir types.Dir) {
	for range only.Once {
		if me.Filepath != "" {
			dir = util.FileDir(me.Filepath)
			break
		}
		fp, sts := me.GetFilepath()
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message: fmt.Sprintf("failed to get filepath for project '%s'", me.Hostname),
			})
			break
		}
		dir = util.FileDir(fp)
	}
	return dir
}

func (me *Project) WriteFile() (sts status.Status) {
	for range only.Once {
		fp, sts := me.GetFilepath()
		if is.Error(sts) {
			break
		}
		jf := NewJsonFile(fp)
		sts = jf.CaptureProject(me)
		if is.Error(sts) {
			break
		}
		sts = jf.Write()
		if is.Error(sts) {
			break
		}
	}
	return sts
}

//func (me *Project) AddNamedStack(stackid gears.StackId) (sts status.Status) {
//	for range only.Once {
//		var stack *gears.NamedStack
//		stack, sts = gears.FindNamedStack(stackid)
//		if is.Error(sts) {
//			break
//		}
//		var sm StackMap
//		sm, sts = stack.GetDefaultServices()
//		if is.Error(sts) {
//			break
//		}
//		for gs, s := range sm {
//			me.Stack[gs] = s
//		}
//		sts = me.WriteFile()
//		if is.Error(sts) {
//			break
//		}
//	}
//	return sts
//}
