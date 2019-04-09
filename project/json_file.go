package project

import (
	"encoding/json"
	"fmt"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var _ util.FilepathHelpUrlGetter = (*JsonFile)(nil)

type HostnameAliases []types.Hostname

type JsonFile struct {
	JsonMeta   JsonMeta               `json:"gearbox"`
	Hostname   types.Hostname         `json:"hostname"`
	Aliases    HostnameAliases        `json:"aliases"`
	ServiceBag gears.ServiceBag       `json:"stack"`
	Stack      service.StackMap       `json:"-"`
	Filepath   types.AbsoluteFilepath `json:"-"`
}

func NewJsonFile(filepath types.AbsoluteFilepath) *JsonFile {
	return &JsonFile{
		Filepath: filepath,
		Stack:    make(service.StackMap, 0),
	}
}

func (me *JsonFile) GetServiceBag() (sb gears.ServiceBag, sts status.Status) {
	sb = make(gears.ServiceBag, len(me.Stack))
	for range only.Once {
		for gs, s := range me.Stack {
			var ps service.Servicer
			ps, sts = s.GetPersistableServiceValue()
			if is.Error(sts) {
				continue
			}
			sb[gs] = ps
		}
	}
	return sb, sts
}

func (me *JsonFile) CaptureProject(project *Project) (sts status.Status) {
	me.Hostname = types.Hostname(project.Hostname)
	me.Aliases = project.Aliases
	me.Stack = project.Stack
	sb, sts := me.GetServiceBag()
	me.ServiceBag = sb
	return sts
}

func (me *JsonFile) Write() (sts status.Status) {
	for range only.Once {
		fp := me.GetFilepath()
		jml := NewJsonMetaLoader(fp)
		sts = jml.Load()
		if status.IsError(sts) {
			break
		}
		if status.IsError(sts) {
			break
		}
		me.JsonMeta = jml.JsonMeta
		b, err := json.MarshalIndent(me, "", "    ")
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to generate JSON for project '%s'", me.Hostname),
			})
			break
		}
		fp = types.AbsoluteFilepath(strings.Replace(string(fp), ".json", "2.json", -1))
		err = ioutil.WriteFile(string(fp), b, os.ModePerm)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to write to '%s'", fp),
				Help:    fmt.Sprintf("ensure you have permissions to write to '%s' and you are not out of disk space", util.FileDir(fp)),
			})
			break
		}
	}
	return
}

type JsonMeta struct {
	Scope       string   `json:"scope"`
	JsonVersion string   `json:"schema"`
	Website     string   `json:"website"`
	Readme      []string `json:"readme"`
}

func (me *JsonFile) GetFilepath() types.AbsoluteFilepath {
	return me.Filepath
}
func (me *JsonFile) GetHelpUrl() string {
	return HelpUrl
}

func (me *JsonFile) Unmarshal(j []byte) (sts status.Status) {
	for range only.Once {
		sts := util.UnmarshalJson(j, me)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("bytes unmarshalled")
	}
	return sts
}

func (me *JsonFile) FixupStack() (sts status.Status) {
	me.Stack = make(service.StackMap, len(me.ServiceBag))
	for gsi, item := range me.ServiceBag {
		var svc *service.Service
		svc, sts = me.FixupStackItem(item, gsi)
		if status.IsError(sts) {
			break
		}
		svc.GearspecId = gsi
		me.Stack[gsi] = service.NewProxyServicer(svc)
	}
	if !status.IsError(sts) {
		me.ServiceBag = nil
		sts = status.Success("stack fixup for '%s' complete", me.Hostname)
	}
	return sts
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
func (me *JsonFile) FixupStackItem(item interface{}, role gearspec.Identifier) (*service.Service, status.Status) {
	var sts status.Status
	ss := service.NewService()
	for range only.Once {
		if svc, ok := item.(string); ok {
			sts = ss.Parse(service.Identifier(svc))
			if status.IsError(sts) {
				break
			}
		} else if gsis, ok := item.([]interface{}); ok {
			services := make(service.Services, len(gsis))
			for i, r := range gsis {
				services[i], sts = me.FixupStackItem(r, role)
				if status.IsError(sts) {
					break
				}
			}
			ss.Services = services
		} else if props, ok := item.(map[string]interface{}); ok {
			var name string
			if name, ok = props["name"].(string); ok {
				ss, sts = me.FixupStackItem(name, role)
			} else {
				sts = status.NewStatus(&status.Args{
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
						sts = status.NewStatus(&status.Args{
							Message:    fmt.Sprintf("Property '%s' if not valid in '%s'", key, role),
							HttpStatus: http.StatusBadRequest,
						})
						break
					}
				}
			}
			if status.IsError(sts) {
				break
			}
		}
	}
	return ss, sts
}

type JsonMetaLoader struct {
	JsonMeta JsonMeta               `json:"gearbox"`
	Filepath types.AbsoluteFilepath `json:"-"`
}

func NewJsonMetaLoader(filepath types.AbsoluteFilepath) *JsonMetaLoader {
	return &JsonMetaLoader{
		Filepath: filepath,
	}
}

func (me *JsonMetaLoader) Load() (sts status.Status) {
	for range only.Once {
		b, err := ioutil.ReadFile(string(me.Filepath))
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to read '%s'", me.Filepath),
				Help: fmt.Sprintf("ensure you have permissions you read '%s'",
					me.Filepath,
				),
			})
			break
		}
		sts = me.UnmarshalMeta(b)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal metadata in '%s'",
					me.Filepath,
				),
			})
			break
		}
	}
	return sts
}

func (me *JsonMetaLoader) UnmarshalMeta(j []byte) (sts status.Status) {
	for range only.Once {
		sts := util.UnmarshalJson(j, me)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("bytes unmarshalled")
	}
	return sts
}

func (me *JsonMetaLoader) GetFilepath() types.AbsoluteFilepath {
	return me.Filepath
}
func (me *JsonMetaLoader) GetHelpUrl() string {
	return HelpUrl
}
