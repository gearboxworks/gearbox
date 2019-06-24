package gears

import (
	"encoding/json"
	"fmt"
	"gearbox/gearspec"
	"gearbox/global"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status/only"

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

type GearRegistrar interface {
	GetName() string
}

type serviceIdsMapGearspecIds map[service.Identifier]gearspec.Identifier

type GearRegistry struct {
	Authorities   Authorities         `json:"authorities"`
	NamedStacks   NamedStacks         `json:"named_stacks"`
	Gearspecs     Gearspecs           `json:"gearspecs"`
	GearOptions   GearOptions         `json:"gear_options"`
	GlobalOptions global.Options      `json:"-"`
	serviceids    service.Identifiers `json:"-"`
	services      Gears               `json:"-"`
	OsBridge      osbridge.OsBridger  `json:"-"`
	serviceIds    serviceIdsMapGearspecIds
	refreshed     bool
}

func NewGearRegistry(ossup osbridge.OsBridger) *GearRegistry {
	return &GearRegistry{
		OsBridge:    ossup,
		Authorities: make(Authorities, 0),
		NamedStacks: make(NamedStacks, 0),
		Gearspecs:   make(Gearspecs, 0),
		GearOptions: make(GearOptions, 0),
	}
}

func (me *GearRegistry) GetAuthorityDomains() (as types.AuthorityDomains, sts Status) {
	as = make(types.AuthorityDomains, 0)
	am := make(map[types.AuthorityDomain]bool, 0)
	for _, gs := range me.Gearspecs {
		_, ok := am[gs.AuthorityDomain]
		if ok {
			continue
		}
		am[gs.AuthorityDomain] = true
		as = append(as, gs.AuthorityDomain)
	}
	return as, nil
}
func (me *GearRegistry) FindAuthority(authority types.AuthorityDomain) (ad types.AuthorityDomain, sts Status) {
	return authority, nil
}

func (me *GearRegistry) FilterByNamedStack(stackid types.StackId) (Gearspecs, Status) {
	return me.Gearspecs.FilterByNamedStack(stackid)
}

func (me *GearRegistry) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *GearRegistry) Initialize() (sts Status) {
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
		for _, gs := range me.Gearspecs {
			sts = gs.Fixup()
			if status.IsError(sts) {
				break
			}
		}
		if is.Error(sts) {
			break
		}
		for _, so := range me.GearOptions {
			sts = so.Fixup(me)
			if is.Error(sts) {
				break
			}
		}
		if is.Error(sts) {
			break
		}
		for _, ns := range me.NamedStacks {
			sts = ns.Fixup(me)
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

func (me *GearRegistry) String() string {
	return string(me.Bytes())
}

func (me *GearRegistry) Bytes() []byte {
	bytes, err := json.Marshal(me)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not unserialize '%s' cache file.\n", JsonFilename))
	}
	return bytes
}

func (me *GearRegistry) Unmarshal(b []byte) (sts Status) {
	err := json.Unmarshal(b, &me)
	if err != nil {
		// @TODO Provide a link to upgrade once we have that established
		sts = status.Wrap(err).
			SetMessage("failed to unmarshal json from '%s': %s", JsonFilename, err.Error()).
			SetAllHelp("Your Gearbox may need to be updated; it may not be compatible with the current JSON schema for '%s' at %s.", JsonFilename, JsonUrl)
	}
	return sts
}

func (me *GearRegistry) FindGearspec(gsid gearspec.Identifier) (gs *gearspec.Gearspec, sts Status) {
	for range only.Once {
		var grgs *Gearspec
		grgs, sts = me.Gearspecs.Find(gsid)
		if is.Error(sts) {
			break
		}
		gs = gearspec.NewGearspecFromGearspecer(grgs)
	}
	return gs, sts
}

func (me *GearRegistry) GetNamedStackIds() (nsids types.StackIds) {
	nsids = make(types.StackIds, len(me.NamedStacks))
	for _, ns := range me.NamedStacks {
		nsids = append(nsids, ns.GetIdentifier())
	}
	nsids.Sort()
	return nsids
}

func (me *GearRegistry) ValidateNamedStackId(stackid types.StackId) (sts Status) {
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

func (me *GearRegistry) FindNamedStack(stackid types.StackId) (ns *NamedStack, sts Status) {
	for range only.Once {
		sts = me.ValidateNamedStackId(stackid)
		if is.Error(sts) {
			break
		}
		nsm := me.NamedStacks.GetMap()
		ns, _ = nsm[stackid]
	}
	return ns, sts
}

func (me *GearRegistry) GetGears() (sm Gears) {
	return me.services
}

func (me *GearRegistry) GetGearIds() (sids service.Identifiers, sts Status) {
	for range only.Once {
		if me.serviceids != nil {
			break
		}
		ss := me.GetGears()
		me.serviceids = make(service.Identifiers, len(ss))
		for i, s := range ss {
			me.serviceids[i] = s.GearId
		}
		me.serviceids.Sort()
	}
	if is.Success(sts) {
		sts = status.Success("got service IDs")
	}
	return me.serviceids, sts
}

func (me *GearRegistry) ValidateGearId(gearid service.Identifier) (sts Status) {
	for range only.Once {
		var ok bool
		for _, sid := range me.serviceids {
			if sid == gearid {
				ok = true
				break
			}
		}
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("service ID '%s' not found", gearid),
				HttpStatus: http.StatusNotFound,
				Help:       fmt.Sprintf("see valid service IDs at %s", JsonUrl),
			})
		} else {
			sts = status.Success("service ID '%s' found", gearid)
		}
	}
	return sts
}

func (me *GearRegistry) FindGear(gearid service.Identifier) (service *Gear, sts Status) {
	var tmp *Gear
	for range only.Once {
		sts = me.ValidateGearId(gearid)
		if is.Error(sts) {
			break
		}
		tmp = NewGear()
		sts = tmp.Parse(gearid)
		if is.Error(sts) {
			break
		}
	}
	if !status.IsError(sts) && tmp != nil {
		service = &Gear{}
		*service = *tmp
	}
	return service, sts

}
