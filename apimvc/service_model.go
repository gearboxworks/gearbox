package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	svc "gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

const ServiceModelType apimodeler.ItemType = "service"

var NilServiceModel = (*ServiceModel)(nil)
var _ apimodeler.Itemer = NilServiceModel

type ServiceModelMap map[gearspec.Identifier]*ServiceModel
type ServiceModels []*ServiceModel

type ServiceModel struct {
	GearspecId  gearspec.Identifier  `json:"gearspec_id,omitempty"`
	ServiceId   types.ServiceId      `json:"service_id,omitempty"`
	ServiceType types.ServiceType    `json:"service_type,omitempty"`
	Orgname     types.Orgname        `json:"orgname,omitempty"`
	Program     types.ProgramName    `json:"program,omitempty"`
	Version     types.Version        `json:"version,omitempty"`
	GearspecIds gearspec.Identifiers `json:"gearspec_ids,omitempty"`
	Gears       *gears.Gears         `json:"-"`
}

func NewModelFromGearsService(ctx *apimodeler.Context, gsvc *gears.Service) (s *ServiceModel, sts status.Status) {
	s = &ServiceModel{
		GearspecId:  gsvc.GearspecId,
		ServiceId:   gsvc.ServiceId,
		Orgname:     gsvc.Orgname,
		Program:     gsvc.Program,
		Version:     gsvc.Version,
		ServiceType: gsvc.ServiceType,
	}
	return s, sts
}

func NewFromServiceService(ctx *apimodeler.Context, gearspecid gearspec.Identifier, ps svc.Servicer) (s *ServiceModel, sts status.Status) {
	var sid svc.Identifier
	var ss *svc.Service
	for range only.Once {
		_ps, ok := ps.(*svc.ServicerProxy)
		if ok {
			ps = _ps.Servicer
		}
		ss, ok = ps.(*svc.Service)
		if ok {
			sid = ss.Identifier
			break
		}
		sid, ok = ps.(svc.Identifier)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("unable to get identifier for service '%s'", gearspecid),
			})
			break
		}
		ss = svc.NewService()
		sts = ss.Parse(sid)
		if is.Error(sts) {
			sts = status.Fail(&status.Args{
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
	gs := gearspec.NewGearspec()
	sts = gs.Parse(gearspecid)
	if is.Success(sts) {
		s = &ServiceModel{
			GearspecId: gs.GetIdentifier(),
			ServiceId:  types.ServiceId(sid),
		}
		for range only.Once {
			if ctx.GetResponseType() != global.ItemResponse {
				break
			}
			if ctx.Controller.GetBasepath() != ServicesBasepath {
				break

			}
			s.Orgname = ss.OrgName
			s.ServiceType = ss.ServiceType
			s.Program = ss.Program
			s.Version = ss.Version.GetIdentifier()
		}
	}
	return s, sts
}

func (me *ServiceModel) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.ServiceId)
}

func (me *ServiceModel) SetId(apimodeler.ItemId) status.Status {
	panic("implement me")
}

func (me *ServiceModel) GetType() apimodeler.ItemType {
	return ServiceModelType
}

func (me *ServiceModel) GetItem() (apimodeler.Itemer, status.Status) {
	return me, nil
}

func (me *ServiceModel) GetItemLinkMap(*apimodeler.Context) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("https://example.com"),
	}, sts
}

func (me *ServiceModel) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.Itemer) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}

func GetServiceModelsFromServiceStackMap(ctx *apimodeler.Context, sm svc.StackMap) (sms ServiceModels, sts status.Status) {
	sms = make(ServiceModels, len(sm))
	i := 0
	for gs, gbs := range sm {
		var s *ServiceModel
		s, sts = NewFromServiceService(ctx, gs, gbs)
		if is.Error(sts) {
			break
		}
		sms[i] = s
		i++
	}
	return sms, sts
}
