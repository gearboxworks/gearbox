package config

import (
	"encoding/json"
	"fmt"
	"gearbox/box"
	"gearbox/help"
	"gearbox/only"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var _ util.FilepathHelpUrlGetter = (*Config)(nil)
var _ Configer = (*Config)(nil)

type Configer interface {
	AddBasedir(*BasedirArgs) (*Basedir, Status)
	DeleteBasedir(types.Nickname) Status
	AddProject(*Project) Status
	Bytes() []byte
	DeleteProject(types.Hostname) Status
	ExpandBasedirPath(types.Nickname, types.RelativePath) (types.AbsoluteDir, Status)
	FindProject(types.Hostname) (*Project, Status)
	FindBasedir(types.Nickname) (*Basedir, Status)
	GetBasedirMap() BasedirMap
	GetNicknameMap() NicknameMap
	GetBasedirNicknames() types.Nicknames
	GetBoxBasedir(types.Nickname) types.AbsoluteDir
	GetCandidates() Candidates
	GetDir() types.AbsoluteDir
	GetFilepath() types.AbsoluteFilepath
	GetHelpUrl() string
	GetBasedir(types.Nickname) (types.AbsoluteDir, Status)
	GetBasedirs() types.AbsoluteDirs
	GetProjectMap() (ProjectMap, Status)
	Initialize() (sts Status)
	Load() Status
	LoadProjects() Status
	LoadProjectsAndWrite() Status
	MakeUniqueBasedirNickname(types.AbsoluteDir) types.Nickname
	MaybeMakeDir(types.AbsoluteDir, os.FileMode) Status
	NicknameExists(types.Nickname) bool
	BasedirExists(types.AbsoluteDir) bool
	Unmarshal(j []byte) Status
	UpdateBasedir(*Basedir) Status
	UpdateProject(*Project) Status
	WriteFile() Status
}

var ProjectRootAddCmd *cobra.Command

type Config struct {
	About         string            `json:"about"`
	LearnMore     string            `json:"learn_more"`
	OsBridge      OsBridger         `json:"-"`
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

func NewConfig(OsBridge OsBridger) Configer {
	c := &Config{
		About:         "This is a Gearbox user configuration file.",
		LearnMore:     "To learn about Gearbox visit https://gearbox.works",
		OsBridge:      OsBridge,
		SchemaVersion: SchemaVersion,
		BasedirMap:    make(BasedirMap, 1),
		ProjectMap:    make(ProjectMap, 0),
		Candidates:    make(Candidates, 0),
		BoxBasedir:    box.Basedir,
	}
	c.BasedirMap[DefaultBasedirNickname] = NewBasedir(
		DefaultBasedirNickname,
		c.OsBridge.GetProjectDir(),
	)
	return c
}

var sanitizer *regexp.Regexp

func init() {
	sanitizer = regexp.MustCompile("[^a-z0-9 ]+")
}

func (me *Config) MakeUniqueBasedirNickname(dir types.AbsoluteDir) (nn types.Nickname) {
	try := strings.ToLower(filepath.Base(string(dir)))
	try = sanitizer.ReplaceAllString(try, "")
	try = strings.Replace(try, " ", "-", -1)
	base := try
	i := 2
	bdm := me.GetBasedirMap()
	for {
		nn = types.Nickname(try)
		bd, ok := bdm[nn]
		if !ok {
			break
		}
		if bd.Basedir == dir {
			break
		}
		try = fmt.Sprintf("%s%d", base, i)
		i++
	}
	return nn
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
		_ = sts.SetHttpStatus(http.StatusCreated)
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
		_ = sts.SetHttpStatus(http.StatusCreated)
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
		_ = sts.SetHttpStatus(http.StatusNoContent)
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
	for nn := range me.BasedirMap {
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

func (me *Config) GetBasedirMap() BasedirMap {
	return me.BasedirMap
}

func (me *Config) GetBasedirs() types.AbsoluteDirs {
	bds := make(types.AbsoluteDirs, len(me.BasedirMap))
	for _, bd := range me.BasedirMap {
		bds = append(bds, bd.Basedir)
	}
	return bds
}

func (me *Config) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}

func (me *Config) GetDir() types.AbsoluteDir {
	return me.OsBridge.GetUserConfigDir()
}

func (me *Config) GetFilepath() types.AbsoluteFilepath {
	fp := filepath.FromSlash(fmt.Sprintf("%s/config.json", me.OsBridge.GetUserConfigDir()))
	return types.AbsoluteFilepath(fp)
}

func (me *Config) WriteFile() (sts Status) {
	for range only.Once {
		j, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marshal config"),
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

func (me *Config) AddBasedir(args *BasedirArgs) (bd *Basedir, sts Status) {
	for range only.Once {
		if args.Nickname != "" {
			sts = status.YourBad("nickname must be empty").
				SetDetail("invalid nickname set as '%s'", args.Nickname)
			break
		}
		if args.Basedir == "" {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid empty directory for '%s'", args.Nickname),
			})
			break
		}
		_bd := Basedir{}
		_bd = Basedir(*args)
		_bd.Nickname = me.MakeUniqueBasedirNickname(_bd.Basedir)
		sts = me.BasedirMap.AddBasedir(me, &_bd)
		if is.Error(sts) {
			if sts.HttpStatus() == http.StatusConflict {
				// Already exists
				break
			}
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("invalid empty directory for '%s'", args.Nickname),
				HttpStatus: http.StatusBadRequest,
			})
			break
		}
		bd = &_bd
		sts = me.WriteFile()
	}
	return bd, sts
}

func (me *Config) GetNamedBasedir(nickname types.Nickname) (bd *Basedir, sts Status) {
	return me.BasedirMap.GetBasedir(nickname)
}

func (me *Config) UpdateBasedir(bd *Basedir) (sts Status) {
	for range only.Once {
		sts = me.BasedirMap.UpdateBasedir(me, bd)
		if is.Error(sts) {
			break
		}
		sts = me.WriteFile()
	}
	return sts
}

func (me *Config) DeleteBasedir(nickname types.Nickname) (sts Status) {
	for range only.Once {
		sts = me.BasedirMap.DeleteBasedir(me, nickname)
		if is.Error(sts) {
			break
		}
		_sts := me.WriteFile()
		if is.Error(_sts) {
			sts = _sts.SetHttpStatus(http.StatusInternalServerError)
			break
		}
	}
	return sts
}

func (me *Config) NicknameExists(nickname types.Nickname) bool {
	return me.BasedirMap.NicknameExists(nickname)
}
func (me *Config) BasedirExists(basedir types.AbsoluteDir) bool {
	return me.BasedirMap.BasedirExists(basedir)
}
