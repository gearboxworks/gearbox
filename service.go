package gearbox

import (
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

type ServiceIds []ServiceId
type ServiceId string

type ServiceMap map[RoleName]*Service

func (me ServiceMap) GetStackNames() StackNames {
	names := util.NewUniqueStrings(len(me))
	for _, s := range me {
		names[s.GetStackName()] = true
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

func (me *Service) GetStackName() (name string) {
	if me.StackRole != nil {
		name = me.StackRole.GetStackName()
	}
	return name
}

func (me *Service) Parse(id string) (status Status) {
	for range only.Once {
		me.Id = ServiceId(id)
		err := me.Identity.Parse(id)
		if err != nil {
			status = NewStatus(&StatusArgs{
				HelpfulError: err.(util.HelpfulError),
				HttpStatus:   http.StatusBadRequest,
			})
		}
		status = NewOkStatus("service id '%s' successfully parsed", id)

	}
	return status
}

func (me *Service) String() string {
	return string(me.Id)
}
