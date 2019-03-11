package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

const ProjectFilename = "gearbox.json"
const ProjectFileHelpUrl = "https://docs.gearbox.works/projects/gearbox.json"

var _ util.FilepathHelpUrlGetter = (*ProjectFile)(nil)

type ProjectFile struct {
	JsonMeta   JsonMeta   `json:"gearbox"`
	Hostname   string     `json:"hostname"`
	Aliases    Aliases    `json:"aliases"`
	StackBag   StackBag   `json:"stack"`
	ServiceMap ServiceMap `json:"-"`
	Filepath   string     `json:"-"`
}

func NewProjectFile(filepath string) *ProjectFile {
	return &ProjectFile{
		Filepath:   filepath,
		ServiceMap: make(ServiceMap, 0),
	}
}

type projectDetails struct {
	Aliases    Aliases    `json:"aliases"`
	ServiceMap ServiceMap `json:"stack"`
	Filepath   string     `json:"filepath"`
}

type JsonMeta struct {
	Scope       string   `json:"scope"`
	JsonVersion string   `json:"json"`
	Version     string   `json:"version"`
	Website     string   `json:"website"`
	Readme      []string `json:"readme"`
}

func (me *ProjectFile) ExportProjectDetails() *projectDetails {
	return &projectDetails{
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

func (me *ProjectFile) Unmarshal(j []byte) (status Status) {
	for range only.Once {
		err := util.UnmarshalJson(j, me)
		if err != nil {
			status = NewStatus(&StatusArgs{
				HelpfulError: err.(util.HelpfulError),
				HttpStatus:   http.StatusInternalServerError,
			})
			break
		}
		status = NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *ProjectFile) FixupStackItem(item interface{}, role RoleName) (*Service, Status) {
	var status Status
	service := NewService(&ServiceArgs{
		StackRole: NewStackRole(),
	})
	for range only.Once {
		if svc, ok := item.(string); ok {
			status = service.Parse(svc)
			if status.IsError() {
				break
			}
			if service.StackRole.NeedsParse() {
				status = service.StackRole.Parse(role)
				if status.IsError() {
					break
				}
			}
			if service.Org == "" {
				service.Org = "gearboxworks"
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
				status = NewStatus(&StatusArgs{
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
						status = NewStatus(&StatusArgs{
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

func (me *ProjectFile) FixupStackMap() (status Status) {
	for role, item := range me.StackBag {
		sr := NewStackRole()
		status = sr.Parse(RoleName(role))
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
		status = NewOkStatus("stack fixup for '%s' complete", me.Hostname)
	}
	return status
}
