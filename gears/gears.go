package gears

import (
	"encoding/json"
	"fmt"
	"gearbox/cache"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"io/ioutil"
	"log"
	"net/http"
)

type Gear interface {
	GetName() string
}

type serviceIdsMapGearspecIds map[service.Identifier]gearspec.Identifier

type Gears struct {
	Authorities       types.Authorities   `json:"authorities"`
	NamedStackIds     types.StackIds      `json:"stacks"`
	StackRoleMap      StackRoleMap        `json:"roles"`
	ServiceOptionsMap RoleServicesMap     `json:"services"`
	GlobalOptions     global.Options      `json:"-"`
	ServiceIds        service.Identifiers `json:"-"`
	ServiceMap        ServiceMap          `json:"-"`
	OsSupport         oss.OsSupporter     `json:"-"`
	serviceIds        serviceIdsMapGearspecIds
	refreshed         bool
}

func NewGears(ossup oss.OsSupporter) *Gears {
	return &Gears{
		OsSupport:         ossup,
		Authorities:       make(types.Authorities, 0),
		NamedStackIds:     make(types.StackIds, 0),
		StackRoleMap:      make(StackRoleMap, 0),
		ServiceOptionsMap: make(RoleServicesMap, 0),
	}
}

func (me *Gears) GetAuthorities() (as types.Authorities, sts status.Status) {
	as = make(types.Authorities, 0)
	am := make(map[types.AuthorityDomain]bool, 0)
	for _, sr := range me.StackRoleMap {
		_, ok := am[sr.Authority]
		if ok {
			continue
		}
		am[sr.Authority] = true
		as = append(as, sr.Authority)
	}
	return as, nil
}
func (me *Gears) FindAuthority(authority types.AuthorityDomain) (ad types.AuthorityDomain, sts status.Status) {
	return authority, nil
}
func (me *Gears) GetStackRoleMap() (StackRoleMap, status.Status) {
	return me.StackRoleMap, nil
}

func (me *Gears) GetNamedStackServiceOptionMap(stackid types.StackId) (rsm RoleServicesMap, sts status.Status) {
	return me.ServiceOptionsMap.FilterForNamedStack(stackid)
}

func (me *Gears) GetNamedStackRoleMap(stackid types.StackId) (StackRoleMap, status.Status) {
	return me.StackRoleMap.FilterByNamedStack(stackid)
}

func (me *Gears) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *Gears) Initialize() (sts status.Status) {
	var b []byte
	if !me.NeedsRefresh() {
		return sts
	}
	for range only.Once {
		cacheDir := me.OsSupport.GetCacheDir()
		store := cache.NewCache(cacheDir)

		store.Disable = me.GlobalOptions.NoCache
		var ok bool
		b, ok, sts = store.Get(CacheKey)
		if ok {
			break
		}
		var sc int
		b, sc, sts = util.HttpRequest(JsonUrl)
		if status.IsError(sts) || sc != http.StatusOK { // @TODO Bundle these as Assets so we will always have some options
			log.Print("Could not download 'gears.json' and no options have previously been stored.")
			filepath := fmt.Sprintf("%s/%s", me.OsSupport.GetAdminRootDir(), JsonFilename)
			var err error
			log.Printf("Loading included '%s'.", filepath)
			b, err = ioutil.ReadFile(filepath)
			if err != nil {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("unable to read '%s'", filepath),
				})
				break
			}
		}
		sts = store.Set(CacheKey, b, "15m")
		if is.Error(sts) {
			log.Printf(sts.Message())
			break
		}
	}
	for range only.Once {
		if is.Error(sts) {
			break
		}
		sts = me.Unmarshal(b)
		if is.Error(sts) {
			break
		}
		for rs, sr := range me.StackRoleMap {
			sts = sr.Fixup(rs)
			if status.IsError(sts) {
				break
			}
		}
		if is.Error(sts) {
			break
		}
		for rs, ro := range me.ServiceOptionsMap {
			sr, ok := me.StackRoleMap[rs]
			if !ok {
				continue // @TODO Log error here and communicate back to home base
			}
			sts = ro.Fixup(sr.GearspecId)
			if is.Error(sts) {
				break
			}
		}
		me.refreshed = true
	}
	if !status.IsError(sts) {
		sts = status.Success("Parent initialized successfully")
	}
	return sts
}

func (me *Gears) String() string {
	return string(me.Bytes())
}

func (me *Gears) Bytes() []byte {
	bytes, err := json.Marshal(me)
	if err != nil {
		log.Fatal("Could not unserialize 'gears.json' cache file.")
	}
	return bytes
}

func (me *Gears) Unmarshal(b []byte) (sts status.Status) {
	err := json.Unmarshal(b, &me)
	if err != nil {
		// @TODO Provide a link to upgrade once we have that established
		sts = status.Wrap(err, &status.Args{
			Message: fmt.Sprintf("failed to unmarshal json from 'gears.json': %s", err.Error()),
			Help: fmt.Sprintf("Your Parent may need to be updated; it may not be compatible with the current JSON schema for 'gears.json' at %s.",
				JsonUrl,
			),
		})
	}
	return sts
}

func (me *Gears) FindGearspec(gsid gearspec.Identifier) (gs *gearspec.Gearspec, sts status.Status) {
	return nil, nil
}

func (me *Gears) GetNamedStackMap() (nsm NamedStackMap, sts status.Status) {
	for range only.Once {
		nsm = make(NamedStackMap, len(me.NamedStackIds))
		for _, nsid := range me.NamedStackIds {
			//			ns := NewNamedStack(me, nsid)
			ns := NewNamedStack(nsid)
			sts = ns.Refresh()
			if is.Error(sts) {
				break
			}
			nsm[nsid] = ns
		}
	}
	return nsm, sts
}

func (me *Gears) GetNamedStackIds() (nsids types.StackIds, sts status.Status) {
	for range only.Once {
		for i, nsid := range me.NamedStackIds {
			me.NamedStackIds[i] = nsid
		}
		me.NamedStackIds.Sort()
	}
	if is.Success(sts) {
		sts = status.Success("got named stack IDs")
	}
	return me.NamedStackIds, sts
}

func (me *Gears) ValidateNamedStackId(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		var ok bool
		for _, nsid := range me.NamedStackIds {
			if nsid == stackid {
				ok = true
				break
			}
		}
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("named stack ID '%s' not found", stackid),
				HttpStatus: http.StatusNotFound,
				Help:       fmt.Sprintf("see valid named stack IDs at %s", JsonUrl),
			})
		} else {
			sts = status.Success("named stack ID '%s' found", stackid)
		}
	}
	return sts
}

func (me *Gears) FindNamedStack(stackid types.StackId) (stack *NamedStack, sts status.Status) {
	var tmp *NamedStack
	for range only.Once {
		sts = me.ValidateNamedStackId(stackid)
		if is.Error(sts) {
			break
		}
		//		tmp = NewNamedStack(me, stackid)
		tmp = NewNamedStack(stackid)
		sts = tmp.Refresh()
		if is.Error(sts) {
			break
		}
	}
	if !status.IsError(sts) && tmp != nil {
		stack = &NamedStack{}
		*stack = *tmp
	}
	return stack, sts
}

func (me *Gears) GetServiceMap() (sm ServiceMap, sts status.Status) {
	for range only.Once {
		if me.ServiceMap != nil {
			break
		}
		me.ServiceMap = make(ServiceMap, 0)
		sids, sts := me.getServiceIdsMapGearspecIds()
		for sid, gsid := range sids {
			if is.Error(sts) {
				break
			}
			s := NewService()
			sts = s.SetIdentifier(sid)
			if is.Error(sts) {
				break
			}
			s.GearspecId = gsid
			me.ServiceMap[sid] = s
		}
	}
	return me.ServiceMap, sts
}

func (me *Gears) getServiceIdsMapGearspecIds() (sids serviceIdsMapGearspecIds, sts status.Status) {
	for range only.Once {
		if me.serviceIds != nil {
			break
		}
		me.serviceIds = make(serviceIdsMapGearspecIds, 0)
		for gsid, so := range me.ServiceOptionsMap {
			for _, s := range so.Services {
				me.serviceIds[s.ServiceId] = gsid
			}
		}
	}
	return me.serviceIds, sts
}

func (me *Gears) GetServiceIds() (sids service.Identifiers, sts status.Status) {
	var ids serviceIdsMapGearspecIds
	for range only.Once {
		if me.ServiceIds != nil {
			break
		}
		ids, sts = me.getServiceIdsMapGearspecIds()
		if is.Error(sts) {
			break
		}
		i := 0
		for sid := range ids {
			me.ServiceIds[i] = sid
			i++
		}
		me.ServiceIds.Sort()
	}
	if is.Success(sts) {
		sts = status.Success("got service IDs")
	}
	return me.ServiceIds, sts
}

func (me *Gears) ValidateServiceId(serviceid service.Identifier) (sts status.Status) {
	for range only.Once {
		var ok bool
		for _, nsid := range me.ServiceIds {
			if nsid == serviceid {
				ok = true
				break
			}
		}
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("service ID '%s' not found", serviceid),
				HttpStatus: http.StatusNotFound,
				Help:       fmt.Sprintf("see valid service IDs at %s", JsonUrl),
			})
		} else {
			sts = status.Success("service ID '%s' found", serviceid)
		}
	}
	return sts
}

func (me *Gears) FindService(serviceid service.Identifier) (service *Service, sts status.Status) {
	var tmp *Service
	for range only.Once {
		sts = me.ValidateServiceId(serviceid)
		if is.Error(sts) {
			break
		}
		tmp = NewService()
		sts = tmp.Parse(serviceid)
		if is.Error(sts) {
			break
		}
	}
	if !status.IsError(sts) && tmp != nil {
		service = &Service{}
		*service = *tmp
	}
	return service, sts

}
