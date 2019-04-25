package config

import (
	"fmt"
	"gearbox/jsonfile"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"path/filepath"
	"strings"
)

type Projects []*Project

type Project struct {
	Hostname types.Hostname     `json:"hostname"`
	Enabled  bool               `json:"enabled"`
	Basedir  types.Nickname     `json:"basedir"`
	Notes    string             `json:"notes"`
	Path     types.RelativePath `json:"path"`
	Config   Configer           `json:"-"`
}

func NewProject(cfg Configer, path types.RelativePath) *Project {
	return &Project{
		Hostname: types.Hostname(path),
		Path:     path,
		Config:   cfg,
	}
}

func (me *Project) GetDir() (dir types.AbsoluteDir, sts Status) {
	for range only.Once {
		dir, sts = me.Config.GetHostBasedir(me.Basedir)
		if status.IsError(sts) {
			break
		}
		dir = types.AbsoluteDir(filepath.FromSlash(fmt.Sprintf("%s/%s", dir, me.Path)))
	}
	return dir, sts
}

func (me *Project) GetFilepath() (fp types.AbsoluteFilepath, sts Status) {
	for range only.Once {
		var bd types.AbsoluteDir
		bd, sts = me.Config.GetHostBasedir(me.Basedir)
		if status.IsError(sts) {
			break
		}
		fp = jsonfile.GetFilepath(bd, me.Path)
	}
	return fp, sts
}

func (me *Project) GetFullpath() (dp types.AbsoluteDir, sts Status) {
	for range only.Once {
		dp, sts = me.Config.ExpandHostBasedirPath(me.Basedir, me.Path)
		if is.Error(sts) {
			break
		}
		dp = types.AbsoluteDir(filepath.FromSlash(string(dp)))
	}
	return dp, sts
}

func (me *Project) GetHostname() types.Hostname {
	hostname := types.Hostname(me.Path)
	if !strings.Contains(string(hostname), ".") {
		hostname = types.Hostname(fmt.Sprintf("%s.local", hostname))
	}
	return types.Hostname(strings.ToLower(string(hostname)))
}

//func ValidateProjectHostname(hostname Hostname, args ...*ValidateArgs) (sts Status) {
//	for range only.Once {
//		var apiHelp string
//		var _args *ValidateArgs
//		if len(args) == 0 {
//			_args = &ValidateArgs{}
//		} else {
//			_args = args[0]
//		}
//		if _args.ApiHelpUrl != "" {
//			apiHelp = fmt.Sprintf("see %s", _args.ApiHelpUrl)
//		}
//
//		if _args.MustNotBeEmpty && hostname == "" {
//			sts = status.Fail(&status.Args{
//				Message:    "project hostname is empty",
//				HttpStatus: http.StatusBadRequest,
//				ApiHelp:    apiHelp,
//			})
//			break
//		}
//		hnExists, sts := _args.ProjectExists(hostname)
//		if is.Error(sts) {
//			sts = status.Wrap(sts, &status.Args{
//				Message: fmt.Sprintf("unabled to load verify project '%s' exists", hostname),
//			})
//			break
//		}
//		if _args.MustExist && !hnExists {
//			sts = status.Fail(&status.Args{
//				Message:    fmt.Sprintf("no project exists with hostname '%s'", hostname),
//				HttpStatus: http.StatusBadRequest,
//				ApiHelp:    apiHelp,
//			})
//			break
//		}
//		if _args.MustNotExist && hnExists {
//			sts = status.Fail(&status.Args{
//				Message:    fmt.Sprintf("project hostname '%s' already exists", hostname),
//				HttpStatus: http.StatusBadRequest,
//				ApiHelp:    apiHelp,
//			})
//			break
//		}
//		sts = status.Success("validated project hostname '%s'", hostname)
//
//	}
//	return sts
//}
