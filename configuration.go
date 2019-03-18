package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/host"
	"gearbox/only"
	"gearbox/stat"
	"gearbox/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var _ util.FilepathHelpUrlGetter = (*Configuration)(nil)
var _ Config = (*Configuration)(nil)

type Config interface {
	Initialize() (status stat.Status)
	GetBasedirNicknames() []string
	GetHostBasedir(string) (string, stat.Status)
	GetHostBasedirs() map[string]string
	Bytes() []byte
	GetDir() string
	GetFilepath() string
	WriteFile() stat.Status
	MaybeMakeDir(string, os.FileMode) stat.Status
	GetHelpUrl() string
	Unmarshal(j []byte) stat.Status
	Load() stat.Status
	LoadProjectsAndWrite() stat.Status
	LoadProjects() stat.Status
	GetProjectMap() ProjectMap
	GetBasedirMap() BasedirMap
	GetBoxBasedir() string
	GetBoxname() string
	GetCandidates() Candidates
}

var ProjectRootAddCmd *cobra.Command

type Configuration struct {
	About         string         `json:"about"`
	LearnMore     string         `json:"learn_more"`
	HostConnector host.Connector `json:"-"`
	SchemaVersion string         `json:"schema_version"`
	BasedirMap    BasedirMap     `json:"basedirs"`
	ProjectMap    ProjectMap     `json:"projects"`
	Candidates    Candidates     `json:"-"`
	BoxBasedir    string         `json:"-"`
	Boxname       string         `json:"-"`
	Gearbox       Gearbox        `json:"-"`
}

func UnmarshalConfig(b []byte) Config {
	c := Configuration{}
	_ = json.Unmarshal(b, &c)
	return &c
}

func NewConfiguration(gb Gearbox) Config {
	c := &Configuration{
		About:         "This is a Gearbox user configuration file.",
		LearnMore:     "To learn about Gearbox visit https://gearbox.works",
		HostConnector: gb.GetHostConnector(),
		SchemaVersion: ConfigSchemaVersion,
		BasedirMap:    make(BasedirMap, 1),
		ProjectMap:    make(ProjectMap, 0),
		Candidates:    make(Candidates, 0),
		BoxBasedir:    BoxBasedir,
		Boxname:       Boxname,
		Gearbox:       gb,
	}
	c.BasedirMap[PrimaryBasedirNickname] = NewBasedir(
		c.HostConnector.GetSuggestedBasedir(),
		&BasedirArgs{
			Nickname: PrimaryBasedirNickname,
		},
	)
	return c
}

func (me *Configuration) Initialize() (status stat.Status) {
	status = me.Load()
	if !status.IsError() {
		status = me.WriteFile()
	}
	return status
}

func (me *Configuration) GetCandidates() Candidates {
	return me.Candidates
}
func (me *Configuration) GetBoxname() string {
	return me.Boxname
}

func (me *Configuration) GetBoxBasedir() string {
	return me.BoxBasedir
}

func (me *Configuration) GetBasedirNicknames() (nns []string) {
	nns = make([]string, len(me.BasedirMap))
	i := 0
	for nn, _ := range me.BasedirMap {
		nns[i] = nn
		i++
	}
	return nns
}

func (me *Configuration) GetHostBasedir(nickname string) (basedir string, status stat.Status) {
	bd, ok := me.BasedirMap[nickname]
	if ok {
		basedir = bd.HostDir
		status = stat.NewOkStatus("hostdir found for nickname '%s'", nickname)
	} else {
		status = stat.NewFailStatus(&stat.Args{
			Error:      stat.IsStatusError,
			Message:    fmt.Sprintf("basedir nickname '%s' is not valid", basedir),
			HttpStatus: http.StatusBadRequest,
			Help: fmt.Sprintf("Add '%s' as a new basedir, or use one of these valid nicknames: %s",
				nickname,
				util.OxfordComma(me.GetBasedirNicknames(), &util.OxfordCommaArgs{
					SingleQuote: true,
					Conjunction: "or",
				}),
			),
		})
	}
	return basedir, status
}

func (me *Configuration) GetBasedirMap() BasedirMap {
	return me.BasedirMap
}

func (me *Configuration) GetHostBasedirs() map[string]string {
	bds := make(map[string]string, len(me.BasedirMap))
	for _, bd := range me.BasedirMap {
		bds[bd.Nickname] = bd.HostDir
	}
	return bds
}

func (me *Configuration) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}

func (me *Configuration) GetDir() string {
	return me.HostConnector.GetUserConfigDir()
}

func (me *Configuration) GetFilepath() string {
	return filepath.FromSlash(fmt.Sprintf("%s/config.json", me.HostConnector.GetUserConfigDir()))
}

func (me *Configuration) WriteFile() (status stat.Status) {
	for range only.Once {
		j, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			status = stat.NewStatus(&stat.Args{
				Message: fmt.Sprintf("unable to marhsal config"),
				Help:    stat.ContactSupportHelp(),
				Error:   err,
			})
			break
		}
		status = me.MaybeMakeDir(me.GetDir(), os.ModePerm)
		if status.IsError() {
			break
		}
		err = ioutil.WriteFile(me.GetFilepath(), j, os.ModePerm)
		if err != nil {
			status = stat.NewStatus(&stat.Args{
				Message: fmt.Sprintf("unable to write to config file '%s'", me.GetFilepath()),
				Help:    fmt.Sprintf("check '%s' for write permissions", filepath.Dir(me.GetFilepath())),
				Error:   err,
			})
			break
		}
		status = stat.NewOkStatus("project config file written")
	}
	return status
}

func (me *Configuration) MaybeMakeDir(dir string, mode os.FileMode) (status stat.Status) {
	for range only.Once {
		err := util.MaybeMakeDir(dir, mode)
		if err == nil {
			status = stat.NewOkStatus("directory '%s' created", dir)
			break
		}
		status = stat.NewStatus(&stat.Args{
			Message: fmt.Sprintf("failed to create directory '%s'", dir),
			Help:    fmt.Sprintf("confirm directory '%s' is readable", filepath.Dir(dir)),
			Error:   err,
		})

	}
	return status
}

func (me *Configuration) ReadBytes() (b []byte, status stat.Status) {
	for range only.Once {
		fp := me.GetFilepath()
		b, status = util.ReadBytes(fp)
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("read %d bytes from file '%s'.", len(b), fp)
	}
	return b, status
}

func (me *Configuration) GetHelpUrl() string {
	return ConfigHelpUrl
}

func (me *Configuration) Unmarshal(j []byte) (status stat.Status) {
	for range only.Once {
		status := util.UnmarshalJson(j, me)
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *Configuration) Load() (status stat.Status) {
	for range only.Once {
		var j []byte
		j, status = me.ReadBytes()
		if status.IsError() {
			break
		}
		if len(j) > 0 {
			status = me.Unmarshal(j)
		}
		if status.IsError() {
			break
		}
		status = me.LoadProjects()
	}
	return status
}

func (me *Configuration) LoadProjectsAndWrite() (status stat.Status) {
	status = me.LoadProjects()
	if !status.IsError() {
		status = me.WriteFile()
	}
	return status
}

func (me *Configuration) GetProjectMap() ProjectMap {
	return me.ProjectMap
}

func (me *Configuration) LoadProjects() (status stat.Status) {
	for range only.Once {
		if len(me.BasedirMap) == 0 {
			status = stat.NewStatus(&stat.Args{
				Message: fmt.Sprintf("no project roots found in %s", me.GetFilepath()),
				CliHelp: fmt.Sprintf("Add with the '%s <dir>' command", ProjectRootAddCmd.CommandPath()),
				ApiHelp: fmt.Sprintf("Add by POSTing JSON to 'add-basedir' resource"),
				Error:   stat.IsStatusError,
			})
			break
		}
		me.Candidates = make(Candidates, 0)
		baseDirs := make([]string, 0)
		for bdnn, bd := range me.BasedirMap {
			baseDirs = append(baseDirs, fmt.Sprintf("'%s'", bd.HostDir)) // For status message
			bd.Nickname = bdnn                                           // In case it is not set, since it is not written to JSON as a property
			var files []os.FileInfo
			if !util.DirExists(bd.HostDir) {
				err := os.Mkdir(bd.HostDir, 0777)
				if err != nil {
					status = stat.NewStatus(&stat.Args{
						Message: fmt.Sprintf("unable to make directory '%s'", bd.HostDir),
						Error:   err,
					})
					break
				}
			}
			files, err := ioutil.ReadDir(bd.HostDir)
			if err != nil {
				status = stat.NewStatus(&stat.Args{
					Message: fmt.Sprintf("unable to read directory %s", bd.HostDir),
					Error:   err,
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
					Path:    file.Name(),
					Gearbox: me.Gearbox,
				})
				if c.IsProject() {
					p := me.ProjectMap.FindProject(bdnn, c.Path)
					if p == nil {
						p = NewProject(me.Gearbox, c.Path)
					} else {
						p.Renew(me.Gearbox, c.Path)
					}
					p.Basedir = bdnn
					me.ProjectMap[p.Hostname] = p
				} else {
					me.Candidates = append(me.Candidates, c)
				}
			}
		}
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

		if status.NotYetFinalized() {
			status = stat.NewOkStatus("projects loaded for basedirs: %s",
				strings.Join(baseDirs, ", "),
			)
		}
	}
	return status
}
