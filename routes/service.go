package routes

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/gearspecid"
	"gearbox/only"
	svc "gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/version"
)

type ServiceMap map[gsid.Identifier]*Service
type Services []*Service

type Service struct {
	ServiceId  types.ServiceId   `json:"service_id"`
	GearspecId gsid.Identifier   `json:"service_type"`
	OrgName    types.OrgName     `json:"org"`
	Program    types.ProgramName `json:"program"`
	Version    *version.Version  `json:"version"`
	Gearbox    gearbox.Gearboxer `json:"-"`
}

func ConvertService(gearspecid gsid.Identifier, ps svc.Servicer) (s *Service) {
	var sts status.Status
	for range only.Once {
		ss, ok := ps.(*svc.Service)
		if ok {
			s = makeServiceFrom(ss.Identifier, ss, gearspecid)
			break
		}
		sid, ok := ps.(svc.Identifier)
		if !ok {
			panic("invalid project servicer")
		}
		ss = svc.NewService()
		sts = ss.Parse(sid)
		if is.Error(sts) {
			panic(fmt.Sprintf("cannot parse identifier for service '%s'", sid))
		}
		s = makeServiceFrom(sid, ss, gearspecid)
	}
	return s
}

func makeServiceFrom(sid svc.Identifier, ss *svc.Service, gsi gsid.Identifier) *Service {
	return &Service{
		GearspecId: gsi,
		ServiceId:  types.ServiceId(sid),
		OrgName:    ss.OrgName,
		Program:    ss.Program,
		Version:    ss.Version,
	}
}

func ConvertServiceMap(sm svc.StackMap) (rsm ServiceMap) {
	rsm = make(ServiceMap, len(sm))
	for gs, gbs := range sm {
		rsm[gs] = ConvertService(gs, gbs)
	}
	return rsm
}
func ConvertServices(sm svc.StackMap) (rss Services) {
	rss = make(Services, len(sm))
	i := 0
	for gs, gbs := range sm {
		rss[i] = ConvertService(gs, gbs)
		i++
	}
	return rss
}
