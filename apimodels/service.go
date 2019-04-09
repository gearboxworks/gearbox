package apimodels

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/only"
	svc "gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/version"
)

type ServiceMap map[gearspec.Identifier]*Service
type Services []*Service

type Service struct {
	ServiceId  types.ServiceId     `json:"service_id"`
	GearspecId gearspec.Identifier `json:"service_type"`
	OrgName    types.OrgName       `json:"org"`
	Program    types.ProgramName   `json:"program"`
	Version    *version.Version    `json:"version"`
}

func ConvertService(gearspecid gearspec.Identifier, ps svc.Servicer) (s *Service, sts status.Status) {
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
	}
	if is.Success(sts) {
		s = makeServiceFrom(sid, ss, gearspecid)
	}
	return s, sts
}

func makeServiceFrom(sid svc.Identifier, ss *svc.Service, gsi gearspec.Identifier) *Service {
	return &Service{
		GearspecId: gsi,
		ServiceId:  types.ServiceId(sid),
		OrgName:    ss.OrgName,
		Program:    ss.Program,
		Version:    ss.Version,
	}
}

func ConvertServiceMap(sm svc.StackMap) (rsm ServiceMap, sts status.Status) {
	rsm = make(ServiceMap, len(sm))
	for gs, gbs := range sm {
		var s *Service
		s, sts = ConvertService(gs, gbs)
		if is.Error(sts) {
			break
		}
		rsm[gs] = s
	}
	return rsm, sts
}
func ConvertServices(sm svc.StackMap) (rss Services, sts status.Status) {
	rss = make(Services, len(sm))
	i := 0
	for gs, gbs := range sm {
		var s *Service
		s, sts = ConvertService(gs, gbs)
		if is.Error(sts) {
			break
		}
		rss[i] = s
		i++
	}
	return rss, sts
}
