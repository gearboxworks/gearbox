package service

import (
	"fmt"
	"gearbox/gear"
	"gearbox/gearspec"
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
	//GetGearspecId() (gearspec.Gearspec, status.Status)
}

type Services []*Service

var _ Servicer = (*Service)(nil)

type Service struct {
	Identifier Identifier `json:"service_id,omitempty"`
	GearspecId gearspec.Identifier
	*gear.Gear
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
	if _args.Gear == nil {
		_args.Gear = gear.NewGear()
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
func (me *Service) GetServiceId() (id Identifier, sts status.Status) {
	if me.OrgName == "" {
		me.OrgName = global.DefaultOrgName
	}
	return Identifier(me.Gear.GetIdentifier()), sts
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
	for range only.Once {
		_gsid, sts := me.GearspecId.GetIdentifier()
		if is.Error(sts) {
			break
		}
		gid := gear.Gear{}
		sts = gid.Parse(gear.Identifier(_gsid))
		if is.Success(sts) && gid.OrgName == global.DefaultOrgName {
			gid.OrgName = ""
		}
		sid = Identifier(gid.String())
	}
	return sid, sts
}

func (me *Service) GetGearspecId() (role gearspec.Identifier, sts status.Status) {
	for range only.Once {
		if me.GearspecId == "" {
			id, _ := me.GetIdentifier()
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("gearspec ID is empty for service '%s'", id),
			})
		}
	}
	return me.GearspecId, nil
}

func (me *Service) GetStackname() (name types.Stackname) {
	panic("implement me")
	return name
}

func (me *Service) GetIdentifier() (id Identifier, sts status.Status) {
	return me.GetServiceId()
}

func (me *Service) String() string {
	_gsid, sts := me.GetIdentifier()
	if is.Error(sts) {
		panic(sts.Message())
	}
	return string(_gsid)
}

func (me *Service) ParseString(serviceid Identifier) (sts status.Status) {
	return me.Parse(Identifier(serviceid))
}

func (me *Service) Parse(serviceid Identifier) (sts status.Status) {
	for range only.Once {
		if me.Gear == nil {
			me.Gear = gear.NewGear()
		}
		sts := me.Gear.Parse(gear.Identifier(serviceid))
		if status.IsError(sts) {
			break
		}
		me.Identifier, sts = me.GetIdentifier()
		if is.Error(sts) {
			break
		}
		sts = status.Success("service serviceid '%s' successfully parsed", serviceid)
	}
	return sts
}
