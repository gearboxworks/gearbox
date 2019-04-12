package gears

import (
	"gearbox/api/global"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

// Example from gears.json/services
var _ = `{
  "gearbox.works/lamp/apache": {
	 "orgname": "gearboxworks",
	 "default": "apache",
	 "options": [
		"apache:2.4"
	 ]
  },
  "gearbox.works/lamp/mysql": {
	 "orgname": "gearboxworks",
	 "default": "mysql",
	 "options": [
		"mysql:5.5",
		"mysql:5.6",
		"mysql:5.7",
		"mysql:8.0"
	 ]
  },
  "gearbox.works/lamp/php": {
	 "orgname": "gearboxworks",
	 "default": "php",
	 "options": [
		"php:5.2",
		"php:5.6",
		"php:7.0",
		"php:7.1",
		"php:7.2"
	 ]
  }
}`

type RoleServicesMap map[gearspec.Identifier]*RoleServices

type RoleServices struct {
	NamedStackId   types.StackId           `json:"-"`
	Orgname        types.Orgname           `json:"orgname,omitempty"`
	ServiceType    types.ServiceType       `json:"service_type,omitempty"`
	Default        types.ServiceId         `json:"default"`
	Shareable      global.ShareableChoices `json:"shareable"`
	ServiceIds     types.ServiceIds        `json:"options,omitempty"`
	DefaultService *Service                `json:"-"`
	Services       Services                `json:"-"`
}

func NewRoleServices(nsid types.StackId) *RoleServices {
	return &RoleServices{
		NamedStackId: nsid,
	}
}

func (me RoleServicesMap) FilterForNamedStack(stackid types.StackId) (nsrm RoleServicesMap, sts status.Status) {
	for range only.Once {
		gsi := gearspec.NewGearspec()
		sts = gsi.SetStackId(stackid)
		if is.Error(sts) {
			break
		}
		stackid = gsi.GetStackId()
		nsrm = make(RoleServicesMap, 0)
		for i, so := range me {
			if so.NamedStackId != stackid {
				continue
			}
			nsrm[i] = so
		}
	}
	return nsrm, sts
}

func (me *RoleServices) Fixup(id gearspec.Identifier) (sts status.Status) {
	for range only.Once {
		gs := gearspec.NewGearspec()
		sts = gs.Parse(gearspec.Identifier(id))
		if is.Error(sts) {
			break
		}
		me.NamedStackId = gs.GetStackId()
		if me.Default != "" {
			me.DefaultService, sts = me.FixupService(me.Default)
			if is.Error(sts) {
				break
			}
		}
		me.Default = ""
		me.Services = make(Services, len(me.ServiceIds))
		for i, sid := range me.ServiceIds {
			var s *Service
			s, sts = me.FixupService(sid)
			if is.Error(sts) {
				break
			}
			s.GearspecId = gs.GetIdentifier()
			me.Services[i] = s
		}
		me.ServiceIds = nil
	}
	return sts
}

func (me *RoleServices) FixupService(serviceId types.ServiceId) (s *Service, sts status.Status) {
	s = NewService()
	sts = s.Parse(serviceId)
	return s, sts
}
