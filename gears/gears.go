package gears

import (
	"encoding/json"
	"fmt"
	"gearbox/cache"
	"gearbox/global"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"log"
	"net/http"
)

type Gear interface {
	GetName() string
}

type Gears struct {
	Authorities       types.Authorities `json:"authorities"`
	NamedStackIds     types.StackIds    `json:"stacks"`
	OsSupport         oss.OsSupporter   `json:"-"`
	StackRoleMap      StackRoleMap      `json:"roles"`
	ServiceOptionsMap ServiceOptionsMap `json:"services"`
	GlobalOptions     global.Options    `json:"-"`
	refreshed         bool
}

func NewGears(ossup oss.OsSupporter) *Gears {
	return &Gears{
		OsSupport:         ossup,
		Authorities:       make(types.Authorities, 0),
		NamedStackIds:     make(types.StackIds, 0),
		StackRoleMap:      make(StackRoleMap, 0),
		ServiceOptionsMap: make(ServiceOptionsMap, 0),
	}
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

func (me *Gears) GetStackRoleMap() (StackRoleMap, status.Status) {
	return me.StackRoleMap, nil
}

func (me *Gears) GetNamedStackServiceOptionMap(stackid types.StackId) (rsm ServiceOptionsMap, sts status.Status) {
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
			log.Fatal("Could not download 'gears.json' and no options have previously been stored.")
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
			Help: fmt.Sprintf("Your Gearbox may need to be updated; it may not be compatible with the current JSON schema for 'gears.json' at %s.",
				JsonUrl,
			),
		})
	}
	return sts
}
