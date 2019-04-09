package gears

import (
	"gearbox/gearid"
	"gearbox/gearspecid"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

// Example from gears.json/services
var _ = `{
  "gearbox.works/lamp/apache": {
	 "org": "gearboxworks",
	 "default": "apache",
	 "options": [
		"apache:2.4"
	 ]
  },
  "gearbox.works/lamp/mysql": {
	 "org": "gearboxworks",
	 "default": "mysql",
	 "options": [
		"mysql:5.5",
		"mysql:5.6",
		"mysql:5.7",
		"mysql:8.0"
	 ]
  },
  "gearbox.works/lamp/php": {
	 "org": "gearboxworks",
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

type RoleServicesMap map[gsid.Identifier]*RoleServices

type RoleServices struct {
	NamedStackId   types.StackId          `json:"-"`
	OrgName        types.OrgName          `json:"org,omitempty"`
	Default        gearid.GearIdentifier  `json:"default"`
	Shareable      ShareableChoices       `json:"shareable"`
	ServiceIds     gearid.GearIdentifiers `json:"choices,omitempty"`
	DefaultService *Service               `json:"-"`
	Services       Services               `json:"-"`
}

func NewRoleServices(nsid types.StackId) *RoleServices {
	return &RoleServices{
		NamedStackId: nsid,
	}
}

func (me RoleServicesMap) FilterForNamedStack(stackid types.StackId) (nsrm RoleServicesMap, sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
		sts = gsi.Parse(gsid.Identifier(stackid))
		if is.Error(sts) {
			break
		}
		stackid = types.StackId(gsi.String())
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

func (me *RoleServices) Fixup(id gsid.Identifier) (sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
		sts = gsi.Parse(gsid.Identifier(id))
		if is.Error(sts) {
			break
		}
		me.NamedStackId = gsi.GetStackId()
		if me.Default != "" {
			me.DefaultService, sts = me.FixupService(me.Default)
			if is.Error(sts) {
				break
			}
			if me.DefaultService.OrgName == "" {
				me.DefaultService.OrgName = me.OrgName
			}
		}
		me.Default = ""
		me.Services = make(Services, len(me.ServiceIds))
		for i, sid := range me.ServiceIds {
			me.Services[i], sts = me.FixupService(sid)
			if is.Error(sts) {
				break
			}
		}
		me.ServiceIds = nil
	}
	return sts
}

func (me *RoleServices) FixupService(serviceId gearid.GearIdentifier) (s *Service, sts status.Status) {
	for range only.Once {
		s = NewService(serviceId)
		if me.DefaultService != nil {
			sts = s.ApplyDefaults(me.DefaultService)
			break
		}
		me.DefaultService = NewService(serviceId)
		_, sts = me.DefaultService.Parse()
		if is.Error(sts) {
			break
		}
		sts = s.ApplyDefaults(me.DefaultService)
	}
	return s, sts
}
