package apimodels

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

const ServiceTypeName apimodeler.ItemType = "service"

var ServiceInstance = (*Service)(nil)
var _ apimodeler.ApiItemer = ServiceInstance

type ServiceMap map[gearspec.Identifier]*Service
type Services []*Service

type Service struct {
	GearspecId  gearspec.Identifier  `json:"gearspec_id,omitempty"`
	ServiceId   types.ServiceId      `json:"service_id,omitempty"`
	ServiceType types.ServiceType    `json:"service_type,omitempty"`
	Orgname     types.Orgname        `json:"orgname,omitempty"`
	Program     types.ProgramName    `json:"program,omitempty"`
	Version     types.Version        `json:"version,omitempty"`
	GearspecIds gearspec.Identifiers `json:"gearspec_ids,omitempty"`
	Gears       *gears.Gears         `json:"-"`
}

func (me *Service) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.ServiceId)
}

func (me *Service) SetId(apimodeler.ItemId) status.Status {
	panic("implement me")
}

func (me *Service) GetType() apimodeler.ItemType {
	return ServiceTypeName
}

func (me *Service) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}

func (me *Service) GetItemLinkMap(*apimodeler.Context) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("https://example.com"),
	}, sts
}

func (me *Service) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.ApiItemer) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}

func NewFromGearsService(ctx *apimodeler.Context, gsvc *gears.Service) (s *Service, sts status.Status) {
	s = &Service{
		GearspecId:  gsvc.GearspecId,
		ServiceId:   gsvc.ServiceId,
		Orgname:     gsvc.Orgname,
		Program:     gsvc.Program,
		Version:     gsvc.Version,
		ServiceType: gsvc.ServiceType,
	}
	return s, sts
}

func GetFromServiceStackMap(ctx *apimodeler.Context, sm svc.StackMap) (rss Services, sts status.Status) {
	rss = make(Services, len(sm))
	i := 0
	for gs, gbs := range sm {
		var s *Service
		s, sts = GetFromServiceService(ctx, gs, gbs)
		if is.Error(sts) {
			break
		}
		rss[i] = s
		i++
	}
	return rss, sts
}

func GetFromServiceService(ctx *apimodeler.Context, gearspecid gearspec.Identifier, ps svc.Servicer) (s *Service, sts status.Status) {
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
		s = &Service{
			GearspecId: gs.GetIdentifier(),
			ServiceId:  types.ServiceId(sid),
		}
		for range only.Once {
			if ctx.GetResponseType() != global.ItemResponse {
				break
			}
			if ctx.Models.Self.GetBasepath() != ServicesBasepath {
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
