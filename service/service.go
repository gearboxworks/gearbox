package service

import (
	"fmt"
	"gearbox/gearid"
	"gearbox/gearspecid"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

type Servicer interface {
	GetPersistableServiceValue() (Servicer, status.Status)
	GetServiceValue() (Servicer, status.Status)
	GetServiceId() (Identifier, status.Status)
	//GetGearspecId() (gears.Gearspec, status.Status)
}

type Services []*Service

var _ Servicer = (*Service)(nil)

type Service struct {
	Identifier Identifier `json:"service_id,omitempty"`
	GearspecId gsid.Identifier
	*gearid.GearId
	Services `json:"services,omitempty"`
}

type Args Service

func NewService(args ...*Args) *Service {
	var _args *Args
	if len(args) == 0 {
		_args = &Args{}
	} else {
		_args = args[0]
	}
	if _args.GearId == nil {
		_args.GearId = gearid.NewGearId()
	}
	svc := Service{}
	svc = Service(*_args)
	return &svc
}

func (me *Service) GetPersistableServiceValue() (Servicer, status.Status) {
	return me.GetPersistableServiceId()
}
func (me *Service) GetServiceValue() (Servicer, status.Status) {
	return me, nil
}
func (me *Service) GetServiceId() (Identifier, status.Status) {
	if me.OrgName == "" {
		me.OrgName = global.DefaultOrgName
	}
	return Identifier(me.GetIdentifier()), nil
}

//
// ServiceId.GetPersistableGearspecId()
//
// Returns a Gearspec without authority if authority is "gearbox.works"
//
// Used to write values to the gearbox.json configuration file
// to keep things simple for the user/reader.
//
func (me *Service) GetPersistableServiceId() (sid Identifier, sts status.Status) {
	gid := gearid.GearId{}
	sts = gid.Parse(gearid.GearIdentifier(me.GetIdentifier()))
	if is.Success(sts) && gid.OrgName == global.DefaultOrgName {
		gid.OrgName = ""
	}
	return Identifier(gid.String()), sts
}

func (me *Service) GetGearspecId() (role gsid.Identifier, sts status.Status) {
	for range only.Once {
		if me.GearspecId == "" {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("gearspec ID is empty for service '%s'", me.GetIdentifier()),
			})
		}
	}
	return me.GearspecId, nil
}

func (me *Service) GetStackname() (name types.Stackname) {
	panic("implement me")
	return name
}

func (me *Service) GetIdentifier() Identifier {
	return Identifier(me.String())
}

func (me *Service) String() string {
	return string(me.GetIdentifier())
}

func (me *Service) ParseString(serviceid Identifier) (sts status.Status) {
	return me.Parse(Identifier(serviceid))
}

func (me *Service) Parse(serviceid Identifier) (sts status.Status) {
	for range only.Once {
		if me.GearId == nil {
			me.GearId = gearid.NewGearId()
		}
		sts := me.GearId.Parse(gearid.GearIdentifier(serviceid))
		if status.IsError(sts) {
			break
		}
		me.Identifier = me.GetIdentifier()
		sts = status.Success("service serviceid '%s' successfully parsed", serviceid)
	}
	return sts
}
