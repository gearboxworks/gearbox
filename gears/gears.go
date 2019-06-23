package gears

import (
	"encoding/json"
	"fmt"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"github.com/gearboxworks/go-osbridge"

	//	"gearbox/os_support"
	"gearbox/service"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-jsoncache"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type Gear interface {
	GetName() string
}

type serviceIdsMapGearspecIds map[service.Identifier]gearspec.Identifier

type Gears struct {
	Authorities    Authorities         `json:"authorities"`
	NamedStacks    NamedStacks         `json:"named_stacks"`
	StackRoles     StackRoles          `json:"stack_roles"`
	ServiceOptions ServiceOptions      `json:"service_options"`
	GlobalOptions  global.Options      `json:"-"`
	serviceids     service.Identifiers `json:"-"`
	services       Services            `json:"-"`
	OsBridge       osbridge.OsBridger  `json:"-"`
	serviceIds     serviceIdsMapGearspecIds
	refreshed      bool
}

func NewGears(ossup osbridge.OsBridger) *Gears {
	return &Gears{
		OsBridge:       ossup,
		Authorities:    make(Authorities, 0),
		NamedStacks:    make(NamedStacks, 0),
		StackRoles:     make(StackRoles, 0),
		ServiceOptions: make(ServiceOptions, 0),
	}
}

func (me *Gears) GetAuthorityDomains() (as types.AuthorityDomains, sts status.Status) {
	as = make(types.AuthorityDomains, 0)
	am := make(map[types.AuthorityDomain]bool, 0)
	for _, sr := range me.StackRoles {
		_, ok := am[sr.AuthorityDomain]
		if ok {
			continue
		}
		am[sr.AuthorityDomain] = true
		as = append(as, sr.AuthorityDomain)
	}
	return as, nil
}
func (me *Gears) FindAuthority(authority types.AuthorityDomain) (ad types.AuthorityDomain, sts status.Status) {
	return authority, nil
}
func (me *Gears) GetStackRoles() (StackRoles, status.Status) {
	return me.StackRoles, nil
}

func (me *Gears) GetNamedServiceOptions(stackid types.StackId) (sos ServiceOptions, sts status.Status) {
	return me.ServiceOptions.FilterForNamedStack(stackid)
}

func (me *Gears) GetNamedStackRoles(stackid types.StackId) (StackRoles, status.Status) {
	return me.StackRoles.FilterByNamedStack(stackid)
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
		cacheDir := me.OsBridge.GetCacheDir()
		store := jsoncache.New(cacheDir)

		store.Disable = me.GlobalOptions.NoCache
		var ok bool
		b, ok, sts = store.Get(CacheKey)
		if ok {
			break
		}
		var sc int
		b, sc, sts = util.HttpRequest(JsonUrl)
		if status.IsError(sts) || sc != http.StatusOK { // @TODO Bundle these as Assets so we will always have some options
			log.Printf("Could not download '%s' and no options have previously been stored.", JsonFilename)
			fp := filepath.FromSlash(fmt.Sprintf("%s/%s", me.OsBridge.GetAdminRootDir(), JsonFilename))
			var err error
			log.Printf("Loading included '%s'.", fp)
			b, err = ioutil.ReadFile(fp)
			if err != nil {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("unable to read '%s'", fp),
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
		for _, sr := range me.StackRoles {
			sts = sr.Fixup()
			if status.IsError(sts) {
				break
			}
		}
		if is.Error(sts) {
			break
		}
		for _, so := range me.ServiceOptions {
			sts = so.Fixup()
			if is.Error(sts) {
				break
			}
		}
		me.refreshed = true
	}
	if !status.IsError(sts) {
		sts = status.Success("Gearbox initialized successfully")
	}
	return sts
}

func (me *Gears) String() string {
	return string(me.Bytes())
}

func (me *Gears) Bytes() []byte {
	bytes, err := json.Marshal(me)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not unserialize '%s' cache file.\n", JsonFilename))
	}
	return bytes
}

func (me *Gears) Unmarshal(b []byte) (sts status.Status) {
	err := json.Unmarshal(b, &me)
	if err != nil {
		// @TODO Provide a link to upgrade once we have that established
		sts = status.Wrap(err).
			SetMessage("failed to unmarshal json from '%s': %s", JsonFilename, err.Error()).
			SetAllHelp("Your Gearbox may need to be updated; it may not be compatible with the current JSON schema for '%s' at %s.", JsonFilename, JsonUrl)
	}
	return sts
}

func (me *Gears) FindGearspec(gsid gearspec.Identifier) (gs *gearspec.Gearspec, sts status.Status) {
	return nil, nil
}

func (me *Gears) GetNamedStackIds() (nsids types.StackIds) {
	nsids = make(types.StackIds, len(me.NamedStacks))
	for _, ns := range me.NamedStacks {
		nsids = append(nsids, ns.GetIdentifier())
	}
	nsids.Sort()
	return nsids
}

func (me *Gears) ValidateNamedStackId(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		var ok bool
		for _, nsid := range me.NamedStacks {
			if nsid.GetIdentifier() == stackid {
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
		sts = tmp.Refresh(me)
		if is.Error(sts) {
			break
		}
	}
	if !status.IsError(sts) && tmp != nil {
		stack = &NamedStack{}
		*stack = *tmp
		stack.Gears = me
	}
	return stack, sts
}

func (me *Gears) GetServices() (sm Services) {
	return me.services
}

func (me *Gears) GetServiceIds() (sids service.Identifiers, sts status.Status) {
	for range only.Once {
		if me.serviceids != nil {
			break
		}
		ss := me.GetServices()
		me.serviceids = make(service.Identifiers, len(ss))
		for i, s := range ss {
			me.serviceids[i] = s.ServiceId
		}
		me.serviceids.Sort()
	}
	if is.Success(sts) {
		sts = status.Success("got service IDs")
	}
	return me.serviceids, sts
}

func (me *Gears) ValidateServiceId(serviceid service.Identifier) (sts status.Status) {
	for range only.Once {
		var ok bool
		for _, nsid := range me.serviceids {
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
