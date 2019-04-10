package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	svc "gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/version"
)

const ServicesName types.RouteName = "services"
const ServicesBasepath types.Basepath = "/services"

type ServiceMap map[gearspec.Identifier]*Service
type Services []*Service

type Service struct {
	GearspecId gearspec.Identifier `json:"gearspec_id,omitempty"`
	ServiceId  types.ServiceId     `json:"service_id,omitempty"`
	OrgName    types.OrgName       `json:"orgname,omitempty"`
	Program    types.ProgramName   `json:"program,omitempty"`
	Version    *version.Version    `json:"version,omitempty"`
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
			s.OrgName = ss.OrgName
			s.Program = ss.Program
			s.Version = ss.Version
		}
	}
	return s, sts
}
