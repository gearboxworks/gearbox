package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/host"
	"gearbox/only"
	"gearbox/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const SchemaVersion = "1.0"
const boxBaseDir = "/home/gearbox/projects"
const boxName = "Gearbox"
const ConfigHelpDocs = "https://docs.gearbox.works/config"

type Config struct {
	About         string         `json:"about"`
	LearnMore     string         `json:"learn_more"`
	HostConnector host.Connector `json:"-"`
	SchemaVersion string         `json:"schema_version"`
	BaseDirs      BaseDirMap     `json:"base_dirs"`
	Projects      ProjectMap     `json:"projects"`
	Candidates    Candidates     `json:"-"`
	BoxBaseDir    string         `json:"-"`
	BoxName       string         `json:"-"`
	Gearbox       *Gearbox       `json:"-"`
}

func UnmarshalConfig(b []byte) *Config {
	c := Config{}
	_ = json.Unmarshal(b, &c)
	return &c
}

func NewConfig(gb *Gearbox) *Config {
	c := &Config{
		About:         "This is a Gearbox user configuration file.",
		LearnMore:     "To learn about Gearbox visit https://gearbox.works",
		HostConnector: gb.HostConnector,
		SchemaVersion: SchemaVersion,
		BaseDirs:      make(BaseDirMap, 1),
		Projects:      make(ProjectMap, 0),
		Candidates:    make(Candidates, 0),
		BoxBaseDir:    boxBaseDir,
		BoxName:       boxName,
		Gearbox:       gb,
	}
	c.BaseDirs[PrimaryBaseDirNickname] = NewBaseDir(
		c.HostConnector.GetSuggestedBaseDir(),
		&BaseDirArgs{
			Nickname: PrimaryBaseDirNickname,
		},
	)
	return c
}

func (me *Config) Initialize() (status Status) {
	status = me.Load()
	if !status.IsError() {
		status = me.Write()
	}
	return status
}

func (me *Config) GetHostBaseDir(nickname string) (basedir string) {
	bd, ok := me.BaseDirs[nickname]
	if ok {
		basedir = bd.HostDir
	}
	return basedir
}

func (me *Config) GetHostBaseDirs() map[string]string {
	bds := make(map[string]string, len(me.BaseDirs))
	for _, bd := range me.BaseDirs {
		bds[bd.Nickname] = bd.HostDir
	}
	return bds
}

func (me *Config) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}

var ProjectRootAddCmd *cobra.Command

func (me *Config) GetDir() string {
	return me.HostConnector.GetUserConfigDir()
}

func (me *Config) GetFilepath() string {
	return fmt.Sprintf("%s/config.json", me.HostConnector.GetUserConfigDir())
}

func (me *Config) Write() (status Status) {
	for range only.Once {
		j, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("unable to marhsal config"),
				Help:       ContactSupportHelp(),
				HttpStatus: http.StatusInternalServerError,
				Error:      err,
			})
			break
		}
		status = me.MaybeMakeDir(me.GetDir(), os.ModePerm)
		if status.IsError() {
			break
		}
		err = ioutil.WriteFile(me.GetFilepath(), j, os.ModePerm)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("unable to write to config file '%s'", me.GetFilepath()),
				Help:       fmt.Sprintf("check '%s' for write permissions", filepath.Dir(me.GetFilepath())),
				HttpStatus: http.StatusInternalServerError,
				Error:      err,
			})
			break
		}
		status = NewOkStatus("project config file written")
	}
	return status
}

func (me *Config) MaybeMakeDir(dir string, mode os.FileMode) (status Status) {
	for range only.Once {
		err := util.MaybeMakeDir(dir, mode)
		if err == nil {
			status = NewOkStatus("directory '%s' created", dir)
			break
		}
		status = NewStatus(&StatusArgs{
			Message:    fmt.Sprintf("failed to create directory '%s'", dir),
			Help:       fmt.Sprintf("confirm directory '%s' is readable", filepath.Dir(dir)),
			HttpStatus: http.StatusInternalServerError,
			Error:      err,
		})

	}
	return status
}

func (me *Config) ReadBytes() (b []byte, status Status) {
	for range only.Once {
		var err error
		fp := me.GetFilepath()
		b, err = ioutil.ReadFile(fp)
		if err != nil && util.ErrorIsFileDoesNotExist(err) {
			err = nil
		}
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("cannot read from '%s' file.", fp),
				Help:       fmt.Sprintf("confirm file '%s' is readable", fp),
				HttpStatus: http.StatusInternalServerError,
				Error:      err,
			})
			break
		}
		status = NewOkStatus("read %d bytes from file '%s'.", len(b), fp)
	}
	return b, status
}

func (me *Config) Unmarshal(j []byte) (status Status) {
	for range only.Once {
		err := json.Unmarshal(j, &me)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message: fmt.Sprintf("unable to load config file '%s'", me.GetFilepath()),
				Help: fmt.Sprintf("ensure config file '%s' is in correct format per %s",
					me.GetFilepath(),
					ConfigHelpDocs,
				),
				HttpStatus: http.StatusInternalServerError,
				Error:      err,
			})
			break
		}
		status = NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *Config) Load() (status Status) {
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

func (me *Config) LoadProjectsAndWrite() (status Status) {
	status = me.LoadProjects()
	if !status.IsError() {
		status = me.Write()
	}
	return status
}

func (me *Config) GetProjectMap() ProjectMap {
	pm := make(ProjectMap, len(me.Projects))
	for hostname, p := range me.Projects {
		pm[hostname] = p
	}
	return pm
}

func (me *Config) LoadProjects() (status Status) {
	for range only.Once {
		if len(me.BaseDirs) == 0 {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("no project roots found in %s", me.GetFilepath()),
				CliHelp:    fmt.Sprintf("Add with the '%s <dir>' command", ProjectRootAddCmd.CommandPath()),
				ApiHelp:    fmt.Sprintf("Add by POSTing JSON to 'add-basedir' resource"),
				HttpStatus: http.StatusInternalServerError,
				Error:      IsStatusError,
			})
			break
		}
		me.Candidates = make(Candidates, 0)
		baseDirs := make([]string, 0)
		for bdnn, bd := range me.BaseDirs {
			baseDirs = append(baseDirs, fmt.Sprintf("'%s'", bd.HostDir)) // For status message
			bd.Nickname = bdnn                                           // In case it is not set, since it is not written to JSON as a property
			var files []os.FileInfo
			if !util.DirExists(bd.HostDir) {
				err := os.Mkdir(bd.HostDir, 0777)
				if err != nil {
					status = NewStatus(&StatusArgs{
						Message:    fmt.Sprintf("unable to make directory '%s'", bd.HostDir),
						HttpStatus: http.StatusInternalServerError,
						Error:      err,
					})
					break
				}
			}
			files, err := ioutil.ReadDir(bd.HostDir)
			if err != nil {
				status = NewStatus(&StatusArgs{
					Message:    fmt.Sprintf("unable to read directory %s", bd.HostDir),
					HttpStatus: http.StatusInternalServerError,
					Error:      err,
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
					BaseDir: bdnn,
					Path:    file.Name(),
					Gearbox: me.Gearbox,
				})
				if c.IsProject() {
					p := me.Projects.FindProject(bdnn, c.Path)
					if p == nil {
						p = NewProject(me.Gearbox, c.Path)
					}
					p.BaseDir = bdnn
					me.Projects[p.Hostname] = p
				} else {
					me.Candidates = append(me.Candidates, c)
				}
			}
		}
		//
		// Remove any old projects that are not located in one of the basedirs
		//
		for k, p := range me.Projects {
			_, ok := me.BaseDirs[p.BaseDir]
			if !ok {
				delete(me.Projects, k)
				continue
			}
		}

		if status.NotYetFinalized() {
			status = NewOkStatus("projects loaded for basedirs: %s",
				strings.Join(baseDirs, ", "),
			)
		}
	}
	return status
}
