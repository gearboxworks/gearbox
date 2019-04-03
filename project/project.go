package project

import (
	"fmt"
	"gearbox/config"
	"gearbox/jsonfile"
	"gearbox/only"
	"gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"net/http"
	"path/filepath"
)

type Map map[types.Hostname]*Project
type Projects []*Project

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

type Project struct {
	loaded bool
	*config.Project
	Filepath types.AbsoluteFilepath `json:"filepath"`
	Aliases  HostnameAliases        `json:"aliases"`
	Stack    service.StackMap       `json:"stack"`
}

func NewProject(configProject *config.Project) *Project {
	return &Project{
		Project: configProject,
	}
}

func (me *Project) Renew(path types.RelativePath) (sts status.Status) {
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
				sts = status.Fail(&status.Args{
					Message:    fmt.Sprintf("unable to get stack service value for gearspecid '%s'", gs),
					Help:       "ensure your gears.json is using the correct roles.", // @TODO improve this help
					HttpStatus: http.StatusBadRequest,
				})
				break
			}
			me.Stack[gs].Servicer = svc
		}
	}
	return sts
}

func (me *Project) Load() (sts status.Status) {
	for range only.Once {
		var fp types.AbsoluteFilepath
		fp, sts = me.GetFilepath()
		if is.Error(sts) {
			break
		}
		var j []byte
		j, sts = util.ReadBytes(fp)
		if is.Error(sts) {
			break
		}
		pf := NewJsonFile(fp)
		if len(j) > 0 {
			sts = pf.Unmarshal(j)
		}
		sts = pf.FixupStack()
		if is.Error(sts) {
			break
		}
		me.loaded = true
	}
	return sts
}

func (me *Project) GetAliases() (aliases HostnameAliases) {
	return me.Aliases
}

func (me *Project) GetServiceMap() (svcmap service.StackMap) {
	return me.Stack
}

func (me *Project) GetFilepath() (fp types.AbsoluteFilepath, sts status.Status) {
	var bd types.AbsoluteDir
	for range only.Once {
		if me.Filepath != "" {
			break
		}
		bd, sts = me.Config.GetHostBasedir(me.Basedir)
		if is.Error(sts) {
			break
		}
		me.Filepath = types.AbsoluteFilepath(filepath.FromSlash(fmt.Sprintf("%s/%s/%s",
			bd,
			me.Path,
			jsonfile.BaseFilename,
		)))
	}
	return me.Filepath, sts
}

func (me *Project) GetProjectDir() (dir types.AbsoluteDir) {
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
