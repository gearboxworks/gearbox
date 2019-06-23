package apimvc

import (
	"fmt"
	"gearbox/apiworks"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

const GearModelType ItemType = "gears"

var NilGearModel = (*GearModel)(nil)
var _ ItemModeler = NilGearModel

type GearModelMap map[gearspec.Identifier]*GearModel
type GearModels []*GearModel

type GearModel struct {
	GearId       service.Identifier  `json:"-"`
	Orgname      types.Orgname       `json:"orgname,omitempty"`
	Program      types.ProgramName   `json:"program,omitempty"`
	Version      types.Version       `json:"version,omitempty"`
	GearRegistry *gears.GearRegistry `json:"-"`
	Model
}

func (me *GearModel) GetAttributeMap() apiworks.AttributeMap {
	panic("implement me")
}

func NewModelFromGear(ctx *Context, gsvc *gears.Gear) (g *GearModel, sts Status) {
	g = &GearModel{
		GearId:  gsvc.GearId,
		Orgname: gsvc.Orgname,
		Program: gsvc.Program,
		Version: gsvc.Version,
	}
	return g, sts
}

func NewModelFromServiceServicer(ctx *Context, ps service.Servicer) (g *GearModel, sts Status) {
	var sid service.Identifier
	var ss *service.Service
	for range only.Once {
		_ps, ok := ps.(*service.ServicerProxy)
		if ok {
			ps = _ps.Servicer
		}
		ss, ok = ps.(*service.Service)
		if ok {
			sid = ss.Identifier
			break
		}
		sid, ok = ps.(service.Identifier)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "unable to get identifier for unknown service",
			})
			break
		}
		ss = service.NewService()
		sts = ss.Parse(sid)
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message: fmt.Sprintf("cannot parse identifier for service '%s'", sid),
			})
			break
		}
		// Make sure it has
		sid, sts = ss.GetServiceId()
		if is.Error(sts) {
			break
		}
	}
	g = &GearModel{
		GearId: service.Identifier(sid),
	}
	for range only.Once {
		// @TODO Add something here
		//if ctx.GetResponseType() != global.ItemResponse {
		//	break
		//}
		//if ctx.Controller.GetBasepath() != GearsBasepath {
		//	break
		//
		//}
		g.Orgname = ss.OrgName
		g.Program = ss.Program
		g.Version = ss.Version.GetIdentifier()
	}
	return g, sts
}

func (me *GearModel) GetId() ItemId {
	return ItemId(me.GearId)
}

func (me *GearModel) SetId(ItemId) Status {
	panic("implement me")
}

func (me *GearModel) GetType() ItemType {
	return GearModelType
}

func GetServiceModelsFromServiceServicerMap(ctx *Context, sm service.ServicerMap) (sms GearModels, sts Status) {
	sms = make(GearModels, len(sm))
	i := 0
	for smgs, gbs := range sm {
		var g *GearModel
		g, sts = NewModelFromServiceServicer(ctx, gbs)
		if is.Error(sts) {
			break
		}
		gs := gearspec.NewGearspec()
		sts = gs.Parse(smgs)
		if is.Error(sts) {
			break
		}
		sms[i] = g
		i++
	}
	return sms, sts
}
