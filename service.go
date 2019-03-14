package gearbox

import (
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

type ServiceIds []ServiceId
type ServiceId string

type ServiceMap map[RoleSpec]*Service

func (me ServiceMap) GetStackNames() StackNames {
	names := util.NewUniqueStrings(len(me))
	for _, s := range me {
		names[string(s.GetStackName())] = true
	}
	return names.ToSlice()
}

type Services []*Service

type Service struct {
	Id ServiceId `json:"service_id,omitempty"`
	*StackRole
	*Identity
	Services `json:"services,omitempty"`
}
type ServiceArgs Service

func NewService(args ...*ServiceArgs) *Service {
	var _args *ServiceArgs
	if len(args) == 0 {
		_args = &ServiceArgs{}
	} else {
		_args = args[0]
	}
	if _args.Identity == nil {
		_args.Identity = NewIdentity()
	}
	svc := Service{}
	svc = Service(*_args)
	return &svc
}

func (me *Service) Assign(serviceId ServiceId, defaultService *Service) {
	me.Parse(serviceId)
	svcId := me.Identity
	for range only.Once {
		if svcId.OrgName == "" {
			if defaultService != nil && defaultService.Identity.OrgName != "" {
				svcId.OrgName = defaultService.Identity.OrgName
			} else {
				svcId.OrgName = me.OrgName
			}
		}
		if defaultService == nil {
			break
		}
		defaultId := defaultService.Identity
		if svcId.Program == "" {
			svcId.Program = defaultId.Program
		}
		if svcId.Type == "" {
			svcId.Type = defaultId.Type
		}
		defaultVersion := defaultId.Version
		if defaultVersion == nil {
			break
		}
		serviceVersion := svcId.Version
		if serviceVersion.Revision == "" {
			serviceVersion.Revision = defaultVersion.Revision
			if serviceVersion.Patch == "" {
				serviceVersion.Patch = defaultVersion.Patch
				if serviceVersion.Minor == "" {
					serviceVersion.Minor = defaultVersion.Minor
					if serviceVersion.Major == "" {
						serviceVersion.Major = defaultVersion.Major
					}
				}
			}
		}
	}
	svcId.raw = string(svcId.GetId())
	me.Id = ServiceId(svcId.raw)
	me.Identity = svcId
}

func (me *Service) GetStackName() (name StackName) {
	if me.StackRole != nil {
		name = me.StackRole.GetStackName()
	}
	return name
}

func (me *Service) Parse(id ServiceId) (status Status) {
	for range only.Once {
		if me.Identity == nil {
			me.Identity = &Identity{}
		}
		err := me.Identity.Parse(string(id))
		if err != nil {
			status = NewStatus(&StatusArgs{
				HelpfulError: err.(util.HelpfulError),
				HttpStatus:   http.StatusBadRequest,
			})
		}
		me.Id = ServiceId(me.GetId())
		status = NewOkStatus("service id '%s' successfully parsed", id)

	}
	return status
}

func (me *Service) String() string {
	return string(me.GetId())
}
