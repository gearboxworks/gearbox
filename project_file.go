package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"gearbox/util"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const ProjectFilename = "gearbox.json"
const ProjectFileHelpUrl = "https://docs.gearbox.works/projects/gearbox.json"

var _ util.FilepathHelpUrlGetter = (*ProjectFile)(nil)

type ProjectFile struct {
	JsonMeta   JsonMeta   `json:"gearbox"`
	Hostname   string     `json:"hostname"`
	Aliases    Aliases    `json:"aliases"`
	ServiceBag ServiceBag `json:"stack"`
	ServiceMap ServiceMap `json:"-"`
	Filepath   string     `json:"-"`
}

func NewProjectFile(filepath string) *ProjectFile {
	return &ProjectFile{
		Filepath:   filepath,
		ServiceMap: make(ServiceMap, 0),
	}
}

func (me *ProjectFile) GetServiceBag() (sb ServiceBag) {
	sb = make(ServiceBag, len(me.ServiceMap))
	for range only.Once {
		for k, s := range me.ServiceMap {
			sb[k] = s.GetFileValue()
		}
	}
	return sb
}

func (me *ProjectFile) CaptureProject(project *Project) {
	me.Hostname = project.Hostname
	me.Aliases = project.Aliases
	me.ServiceMap = project.ServiceMap
	me.ServiceBag = me.GetServiceBag()
}

func (me *ProjectFile) WriteFile() (status stat.Status) {
	for range only.Once {
		fp := me.GetFilepath()
		jml := NewJsonMetaLoader(fp)
		status = jml.Load()
		if status.IsError() {
			break
		}
		me.JsonMeta = jml.JsonMeta
		b, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("unable to generate JSON for project '%s'", me.Hostname),
			})
			break
		}
		fp = strings.Replace(fp, ".json", "2.json", -1)
		err = ioutil.WriteFile(fp, b, os.ModePerm)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("unable to write to '%s'", fp),
				Help:    fmt.Sprintf("ensure you have permissions to write to '%s' and you are not out of disk space", filepath.Dir(fp)),
			})
			break
		}
	}
	return
}

type ProjectDetails struct {
	Filepath   string     `json:"filepath"`
	Aliases    Aliases    `json:"aliases"`
	ServiceMap ServiceMap `json:"stack"`
}

type JsonMeta struct {
	Scope       string   `json:"scope"`
	JsonVersion string   `json:"schema"`
	Website     string   `json:"website"`
	Readme      []string `json:"readme"`
}

func (me *ProjectFile) ExportProjectDetails() *ProjectDetails {
	return &ProjectDetails{
		Aliases:    me.Aliases,
		ServiceMap: me.ServiceMap,
		Filepath:   me.Filepath,
	}
}

func (me *ProjectFile) GetFilepath() string {
	return me.Filepath
}
func (me *ProjectFile) GetHelpUrl() string {
	return ProjectFileHelpUrl
}

func (me *ProjectFile) Unmarshal(j []byte) (status stat.Status) {
	for range only.Once {
		status := util.UnmarshalJson(j, me)
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *ProjectFile) FixupStack() (status stat.Status) {
	me.ServiceMap = make(ServiceMap, len(me.ServiceBag))
	for role, item := range me.ServiceBag {
		sr := NewStackRole()
		status = sr.Parse(RoleSpec(role))
		if status.IsError() {
			break
		}
		var service *Service
		service, status = me.FixupStackItem(item, role)
		if status.IsError() {
			break
		}
		service.StackRole = sr
		me.ServiceMap[role] = service
	}
	if !status.IsError() {
		me.ServiceBag = nil
		status = stat.NewOkStatus("stack fixup for '%s' complete", me.Hostname)
	}
	return status
}

//
// This processes stack items (services) to allow a service to be specified as any of:
//
//		1. A service ID string
//		2. An array of service ID strings
//		3. A service object
//		4. An array of service objects
//
// Stacks are loaded as a map[string]interface{} to enable this type of processing.
//
func (me *ProjectFile) FixupStackItem(item interface{}, role RoleSpec) (*Service, stat.Status) {
	var status stat.Status
	service := NewService(&ServiceArgs{
		StackRole: NewStackRole(),
	})
	for range only.Once {
		if svc, ok := item.(string); ok {
			status = service.Parse(ServiceId(svc))
			if status.IsError() {
				break
			}
			if service.StackRole.NeedsParse() {
				status = service.StackRole.Parse(role)
				if status.IsError() {
					break
				}
			}
			if service.OrgName == "" {
				service.OrgName = "gearboxworks"
			}
		} else if roles, ok := item.([]interface{}); ok {
			services := make(Services, len(roles))
			for i, r := range roles {
				services[i], status = me.FixupStackItem(r, role)
				if status.IsError() {
					break
				}
			}
			service.Services = services
		} else if props, ok := item.(map[string]interface{}); ok {
			var name string
			if name, ok = props["name"].(string); ok {
				service, status = me.FixupStackItem(name, role)
			} else {
				status = stat.NewStatus(&stat.Args{
					Message:    fmt.Sprintf("Property 'name' if not a string in '%s'", role),
					HttpStatus: http.StatusBadRequest,
				})
				break
			}
			for key, value := range props {
				switch key {
				case "---":
					_, ok := value.(string)
					if ok {
						// Capture the value
					} else {
						status = stat.NewStatus(&stat.Args{
							Message:    fmt.Sprintf("Property '%s' if not valid in '%s'", key, role),
							HttpStatus: http.StatusBadRequest,
						})
						break
					}
				}
			}
			if status.IsError() {
				break
			}
		}
	}
	return service, status
}

type JsonMetaLoader struct {
	JsonMeta JsonMeta `json:"gearbox"`
	Filepath string   `json:"-"`
}

func NewJsonMetaLoader(filepath string) *JsonMetaLoader {
	return &JsonMetaLoader{
		Filepath: filepath,
	}
}

func (me *JsonMetaLoader) Load() (status stat.Status) {
	for range only.Once {
		b, err := ioutil.ReadFile(me.Filepath)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("unable to read '%s'", me.Filepath),
				Help: fmt.Sprintf("ensure you have permissions you read '%s'",
					me.Filepath,
				),
			})
			break
		}
		status = me.UnmarshalMeta(b)
		if status.IsError() {
			status.Status = status
			status.Message = fmt.Sprintf("unable to unmarshal metadata in '%s'",
				me.Filepath,
			)
			break
		}
	}
	return status
}

func (me *JsonMetaLoader) UnmarshalMeta(j []byte) (status stat.Status) {
	for range only.Once {
		status := util.UnmarshalJson(j, me)
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *JsonMetaLoader) GetFilepath() string {
	return me.Filepath
}
func (me *JsonMetaLoader) GetHelpUrl() string {
	return ProjectFileHelpUrl
}
