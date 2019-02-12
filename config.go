package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/host"
	"gearbox/only"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

const SchemaVersion = "1.0"
const vmProjectRoot = "/home/gearbox/projects"

type Config struct {
	About         string         `json:"about"`
	LearnMore     string         `json:"learn_more"`
	HostConnector host.Connector `json:"-"`
	SchemaVersion string         `json:"schema_version"`
	ProjectRoots  ProjectRoots   `json:"project_roots"`
	Projects      Projects       `json:"projects"`
	Candidates    Candidates     `json:"-"`
	VmProjectRoot string         `json:"-"`
}

func NewConfig(hc host.Connector) *Config {
	c := &Config{
		About:         "This is a Gearbox user configuration file.",
		LearnMore:     "To learn about Gearbox visit https://gearbox.works",
		HostConnector: hc,
		SchemaVersion: SchemaVersion,
		ProjectRoots:  make(ProjectRoots, 1),
		Projects:      make(Projects, 0),
		Candidates:    make(Candidates, 0),
		VmProjectRoot: vmProjectRoot,
	}
	c.ProjectRoots[0] = &ProjectRoot{
		HostDir: c.HostConnector.GetSuggestedProjectRoot(),
		VmDir:   vmProjectRoot,
	}
	return c
}

func (me *Config) Initialize() {
	for range only.Once {
		file := me.GetFilepath()
		_, err := os.Stat(file)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			log.Fatal(err.Error())
		}
	}
	me.Load()
	me.Write()
}

var ProjectRootAddCmd *cobra.Command

func (me *Config) GetFilepath() string {
	return fmt.Sprintf("%s/config.json", me.HostConnector.GetUserConfigDir())
}

func (me *Config) Write() {
	j, err := json.MarshalIndent(me, "", "    ")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = ioutil.WriteFile(me.GetFilepath(), j, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (me *Config) Load() {
	for range only.Once {
		j, err := ioutil.ReadFile(me.GetFilepath())
		if err != nil {
			log.Fatal(err.Error())
		}
		err = json.Unmarshal(j, &me)
		if err != nil {
			log.Fatal(err.Error())
		}
		me.LoadProjects()
	}
}

func (me *Config) LoadProjectsAndWrite() {
	me.LoadProjects()
	me.Write()
}

func (me *Config) GetProjectMap() ProjectMap {
	pm := make(ProjectMap, len(me.Projects))
	for _, p := range me.Projects {
		pm[p.Name] = p
	}
	return pm
}

func (me *Config) LoadProjects() {
	var err error
	if len(me.ProjectRoots) == 0 {
		log.Fatal(fmt.Sprintf("No project roots found in %s. Add with the '%s <dir>' command.",
			me.GetFilepath(),
			ProjectRootAddCmd.CommandPath(),
		))
	}
	projectMap := me.GetProjectMap()
	me.Projects = make(Projects, 0)
	me.Candidates = make(Candidates, 0)
	for index, pr := range me.ProjectRoots {
		group := index + 1
		var files []os.FileInfo
		files, err = ioutil.ReadDir(pr.HostDir)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			if file.Name()[0] == '.' {
				continue
			}
			c := &Candidate{
				Root:  &pr.HostDir,
				Path:  file.Name(),
				Group: group,
			}
			if c.IsProject() {
				p, ok := projectMap[c.Path]
				if ok {
					p.Root = c.Root
					p.Hostname = p.MakeHostname()
				} else {
					p = NewProject(c.Path, c.Root)
				}
				p.Group = group
				me.Projects = append(me.Projects, p)
			} else {
				me.Candidates = append(me.Candidates, c)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}
