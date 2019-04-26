package config

import (
	"encoding/json"
	"fmt"
	"gearbox/box"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var _ util.FilepathHelpUrlGetter = (*Config)(nil)
var _ Configer = (*Config)(nil)

type Configer interface {
	AddBasedir(types.AbsoluteDir, ...types.Nickname) Status
	AddProject(*Project) Status
	Bytes() []byte
	DeleteProject(types.Hostname) Status
	ExpandBasedirPath(types.Nickname, types.RelativePath) (types.AbsoluteDir, Status)
	FindProject(types.Hostname) (*Project, Status)
	FindBasedir(types.Nickname) (*Basedir, Status)
	GetBasedirMap() BasedirMap
	GetBasedirNicknames() types.Nicknames
	GetBoxBasedir(types.Nickname) types.AbsoluteDir
	GetCandidates() Candidates
	GetDir() types.AbsoluteDir
	GetFilepath() types.AbsoluteFilepath
	GetHelpUrl() string
	GetBasedir(types.Nickname) (types.AbsoluteDir, Status)
	GetBasedirs() map[types.Nickname]types.AbsoluteDir
	GetProjectMap() (ProjectMap, Status)
	Initialize() (sts Status)
	Load() Status
	LoadProjects() Status
	LoadProjectsAndWrite() Status
	MaybeMakeDir(types.AbsoluteDir, os.FileMode) Status
	NamedBasedirExists(types.Nickname) bool
	Unmarshal(j []byte) Status
	UpdateBasedir(types.Nickname, types.AbsoluteDir) Status
	UpdateProject(*Project) Status
	WriteFile() Status
}

var ProjectRootAddCmd *cobra.Command

type Config struct {
	About         string            `json:"about"`
	LearnMore     string            `json:"learn_more"`
	OsSupport     oss.OsSupporter   `json:"-"`
	SchemaVersion string            `json:"schema_version"`
	BasedirMap    BasedirMap        `json:"basedirs"`
	ProjectMap    ProjectMap        `json:"projects"`
	Candidates    Candidates        `json:"-"`
	BoxBasedir    types.AbsoluteDir `json:"-"`
	Boxname       string            `json:"-"`
}

func UnmarshalConfig(b []byte) Configer {
	c := Config{}
	_ = json.Unmarshal(b, &c)
	return &c
}

func NewConfig(OsSupport oss.OsSupporter) Configer {
	c := &Config{
		About:         "This is a Gearbox user configuration file.",
		LearnMore:     "To learn about Gearbox visit https://gearbox.works",
		OsSupport:     OsSupport,
		SchemaVersion: SchemaVersion,
		BasedirMap:    make(BasedirMap, 1),
		ProjectMap:    make(ProjectMap, 0),
		Candidates:    make(Candidates, 0),
		BoxBasedir:    box.Basedir,
	}
	c.BasedirMap[DefaultBasedirNickname] = NewBasedir(
		DefaultBasedirNickname,
		c.OsSupport.GetSuggestedBasedir(),
	)
	return c
}

func (me *Config) AddProject(p *Project) (sts Status) {
	for range only.Once {
		pm, sts := me.GetProjectMap()
		if status.IsError(sts) {
			break
		}
		_, exists := pm[p.Hostname]
		if exists {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("project '%s' created", p.Hostname),
				HttpStatus: http.StatusConflict,
				Help:       fmt.Sprintf("you cannot create a project with hostname '%s' as it already exists", p.Hostname),
			})
		}
		pm[p.Hostname] = p
		sts = me.WriteFile()
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' created", p.Hostname)
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func (me *Config) UpdateProject(p *Project) (sts Status) {
	for range only.Once {
		pm, sts := me.GetProjectMap()
		if status.IsError(sts) {
			break
		}
		pm[p.Hostname] = p
		sts = me.WriteFile()
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' created", p.Hostname)
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func (me *Config) DeleteProject(hostname types.Hostname) (sts Status) {
	for range only.Once {
		pm, sts := me.GetProjectMap()
		if status.IsError(sts) {
			break
		}
		_, exists := pm[hostname]
		delete(pm, hostname)
		sts = me.WriteFile()
		if status.IsError(sts) {
			break
		}
		if exists {
			sts = status.Success("project '%s' deleted", hostname)
		} else {
			sts = status.Success("project '%s' not found", hostname)
		}
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts
}

func (me *Config) Initialize() (sts Status) {
	sts = me.Load()
	if status.IsError(sts) {
		sts = me.WriteFile()
	}
	return sts
}

func (me *Config) GetCandidates() Candidates {
	return me.Candidates
}

func (me *Config) GetBoxBasedir(nickname types.Nickname) types.AbsoluteDir {
	return types.AbsoluteDir(
		strings.Replace(
			BoxBasedirTemplate,
			NicknameTemplateVar,
			string(nickname),
			-1,
		),
	)
}

func (me *Config) GetBasedirNicknames() (nns types.Nicknames) {
	nns = make(types.Nicknames, len(me.BasedirMap))
	i := 0
	for nn, _ := range me.BasedirMap {
		nns[i] = nn
		i++
	}
	return nns
}

func (me *Config) GetBasedir(nickname types.Nickname) (basedir types.AbsoluteDir, sts Status) {
	bd, ok := me.BasedirMap[nickname]
	if ok {
		basedir = bd.Basedir
		sts = status.Success("hostdir found for nickname '%s'", nickname)
	} else {
		sts = status.Fail(&status.Args{
			Message:    fmt.Sprintf("basedir nickname '%s' is not valid", basedir),
			HttpStatus: http.StatusBadRequest,
			Help: fmt.Sprintf("Add '%s' as a new basedir, or use one of these valid nicknames: %s",
				nickname,
				util.OxfordComma(me.GetBasedirNicknames().Strings(), &util.OxfordCommaArgs{
					SingleQuote: true,
					Conjunction: "or",
				}),
			),
		})
	}
	return basedir, sts
}

func (me *Config) GetBasedirMap() BasedirMap {
	return me.BasedirMap
}

func (me *Config) FindBasedir(nickname types.Nickname) (bd *Basedir, sts Status) {
	bd, ok := me.BasedirMap[nickname]
	if !ok {
		sts = status.Fail(&status.Args{
			Message:    fmt.Sprintf("basedir '%s' not found", nickname),
			HttpStatus: http.StatusNotFound,
		})
	}
	return bd, sts
}

func (me *Config) GetBasedirs() map[types.Nickname]types.AbsoluteDir {
	bds := make(map[types.Nickname]types.AbsoluteDir, len(me.BasedirMap))
	for _, bd := range me.BasedirMap {
		bds[bd.Nickname] = bd.Basedir
	}
	return bds
}

func (me *Config) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}

func (me *Config) GetDir() types.AbsoluteDir {
	return me.OsSupport.GetUserConfigDir()
}

func (me *Config) GetFilepath() types.AbsoluteFilepath {
	fp := filepath.FromSlash(fmt.Sprintf("%s/config.json", me.OsSupport.GetUserConfigDir()))
	return types.AbsoluteFilepath(fp)
}

func (me *Config) WriteFile() (sts Status) {
	for range only.Once {
		j, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marhsal config"),
				Help:    help.ContactSupportHelp(),
			})
			break
		}
		sts = me.MaybeMakeDir(me.GetDir(), os.ModePerm)
		if status.IsError(sts) {
			break
		}
		err = ioutil.WriteFile(string(me.GetFilepath()), j, os.ModePerm)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to write to config file '%s'", me.GetFilepath()),
				Help:    fmt.Sprintf("check '%s' for write permissions", util.FileDir(me.GetFilepath())),
			})
			break
		}
		sts = status.Success("project config file written")
	}
	return sts
}

func (me *Config) MaybeMakeDir(dir types.AbsoluteDir, mode os.FileMode) (sts Status) {
	for range only.Once {
		err := util.MaybeMakeDir(dir, mode)
		if err == nil {
			sts = status.Success("directory '%s' created", dir)
			break
		}
		sts = status.Wrap(err, &status.Args{
			Message: fmt.Sprintf("failed to create directory '%s'", dir),
			Help:    fmt.Sprintf("confirm directory '%s' is readable", util.ParentDir(dir)),
		})

	}
	return sts
}

func (me *Config) ReadBytes() (b []byte, sts Status) {
	for range only.Once {
		fp := me.GetFilepath()
		b, sts = util.ReadBytes(fp)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("read %d bytes from file '%s'.", len(b), fp)
	}
	return b, sts
}

func (me *Config) GetHelpUrl() string {
	return HelpUrl
}

func (me *Config) Unmarshal(j []byte) (sts Status) {
	for range only.Once {
		sts := util.UnmarshalJson(j, me)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("bytes unmarshalled")
	}
	return sts
}

func (me *Config) Load() (sts Status) {
	for range only.Once {
		var j []byte
		j, sts = me.ReadBytes()
		if status.IsError(sts) {
			break
		}
		if len(j) > 0 {
			sts = me.Unmarshal(j)
		}
		if status.IsError(sts) {
			break
		}
		sts = me.LoadProjects()
	}
	return sts
}

func (me *Config) LoadProjectsAndWrite() (sts Status) {
	sts = me.LoadProjects()
	if !status.IsError(sts) {
		sts = me.WriteFile()
	}
	return sts
}

func (me *Config) GetProjectMap() (pm ProjectMap, sts Status) {
	for range only.Once {
		if me.ProjectMap != nil {
			break
		}
		sts = me.LoadProjects()
	}
	return me.ProjectMap, sts
}

func (me *Config) FindProject(hostname types.Hostname) (*Project, Status) {
	return me.ProjectMap.FindProject(hostname)
}

func (me *Config) LoadProjects() (sts Status) {
	for range only.Once {
		if len(me.BasedirMap) == 0 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("no project roots found in %s", me.GetFilepath()),
				CliHelp: fmt.Sprintf("Add with the '%s <dir>' command", ProjectRootAddCmd.CommandPath()),
				ApiHelp: fmt.Sprintf("Add by POSTing JSON to 'add-basedir' routes"),
			})
			break
		}
		me.Candidates = make(Candidates, 0)
		baseDirs := make([]string, 0)
		for bdnn, bd := range me.BasedirMap {
			baseDirs = append(baseDirs, fmt.Sprintf("'%s'", bd.Basedir)) // For status message
			bd.Nickname = bdnn                                           // In case it is not set, since it is not written to JSON as a property
			var files []os.FileInfo
			if !util.DirExists(bd.Basedir) {
				err := os.Mkdir(string(bd.Basedir), 0777)
				if err != nil {
					sts = status.Wrap(err, &status.Args{
						Message: fmt.Sprintf("unable to make directory '%s'", bd.Basedir),
					})
					break
				}
			}
			files, err := ioutil.ReadDir(string(bd.Basedir))
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("unable to read directory %s", bd.Basedir),
				})
				break
			}
			for _, file := range files {
				if !file.IsDir() {
					continue
				}
				if file.Name()[0] == '.' {
					continue
				}
				c := NewCandidate(&CandidateArgs{
					Config:  me,
					Basedir: bdnn,
					Path:    types.RelativePath(file.Name()),
				})
				if !c.IsProject() {
					me.Candidates = append(me.Candidates, c)
				} else {
					p, _ := me.ProjectMap.FindProjectByPath(bdnn, c.Path)
					if p == nil {
						p = NewProject(me, c.Path)
					}
					if is.Error(sts) {
						break
					}
					p.Basedir = bdnn
					p.Config = me
					me.ProjectMap[p.Hostname] = p
				}
			}
			if is.Error(sts) {
				break
			}
		}
		if !status.IsError(sts) {
			//
			// Remove any old projects that are not located in one of the basedirs
			//
			for k, p := range me.ProjectMap {
				_, ok := me.BasedirMap[p.Basedir]
				if !ok {
					delete(me.ProjectMap, k)
					continue
				}
			}
			sts = status.Success("projects loaded for basedirs: %s",
				strings.Join(baseDirs, ", "),
			)
		}
	}
	return sts
}

func (me *Config) ExpandBasedirPath(nickname types.Nickname, path types.RelativePath) (fp types.AbsoluteDir, sts Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
		})
		if is.Error(sts) {
			break
		}
		bd, sts := me.GetBasedir(nickname)
		if is.Error(sts) {
			break
		}
		fp = types.AbsoluteDir(filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path)))
	}
	return fp, sts
}

func (me *Config) AddBasedir(dir types.AbsoluteDir, nickname ...types.Nickname) (sts Status) {
	for range only.Once {
		var nn types.Nickname
		if len(nickname) > 0 {
			nn = nickname[0]
		}
		if dir == "" {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid empty directory for '%s'", nn),
			})
			break
		}
		bd := NewBasedir(nn, dir)
		sts = me.BasedirMap.AddBasedir(bd)
	}
	return sts
}

func (me *Config) GetNamedBasedir(nickname types.Nickname) (bd *Basedir, sts Status) {
	return me.BasedirMap.GetNamedBasedir(nickname)
}

func (me *Config) UpdateBasedir(nickname types.Nickname, dir types.AbsoluteDir) (sts Status) {
	return me.BasedirMap.UpdateBasedir(nickname, dir)
}

func (me *Config) DeleteNamedBasedir(nickname types.Nickname) (sts Status) {
	return me.BasedirMap.DeleteNamedBasedir(nickname)
}

func (me *Config) NamedBasedirExists(nickname types.Nickname) bool {
	return me.BasedirMap.NamedBasedirExists(nickname)
}
