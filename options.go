package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/cache"
	"gearbox/only"
	"gearbox/stat"
	"gearbox/util"
	"log"
	"net/http"
	"strings"
)

const GearsJsonUrl = RepoRawBaseUrl + "/master/assets/gears.json"
const GearsKey = "gears"

type Gears struct {
	Gearbox         Gearbox         `json:"-"`
	Authorities     Authorities     `json:"authorities"`
	StackNames      StackNames      `json:"stacks"`
	RoleMap         RoleMap         `json:"roles"`
	RoleServicesMap RoleServicesMap `json:"role_services"`
	refreshed       bool
}

func NewGears(gb Gearbox) *Gears {
	o := Gears{
		Gearbox: gb,
	}
	return &o
}

//type StackServiceMap map[StackName]*StackService
//type StackService struct {
//	Roles RoleServicesMap `json:"roles"`
//}

type ShareableChoices string

const (
	NotShareable    ShareableChoices = "no"
	InStackSharable ShareableChoices = "instack"
	YesShareable    ShareableChoices = "yes"
)

type RoleServicesMap map[RoleSpec]*RoleService

func (me RoleServicesMap) GetStackServicesMap(stackName StackName) (om RoleServicesMap, status stat.Status) {
	for range only.Once {
		om = make(RoleServicesMap, 0)
		for rs, o := range me {
			stackName, status = GetFullStackName(stackName)
			if status.IsError() {
				break
			}
			if !strings.HasPrefix(string(rs), string(stackName)) {
				continue
			}
			om[rs] = o
		}
	}
	return om, status
}

type RoleService struct {
	*StackRole
	OrgName        OrgName          `json:"org,omitempty"`
	Default        ServiceId        `json:"default,omitempty"`
	Shareable      ShareableChoices `json:"shareable,omitempty"`
	ServiceIds     ServiceIds       `json:"options,omitempty"`
	DefaultService *Service         `json:"-"`
	ServiceOptions Services         `json:"-"`
}

func (me *Gears) GetStackRoleMap(stackName StackName) RoleMap {
	srs := make(RoleMap, 0)
	for i, r := range me.RoleMap {
		if r.GetStackName() != stackName {
			continue
		}
		srs[i] = r
	}
	return srs
}

func (me *RoleService) Fixup(stackRole *StackRole) {
	if me.Default != "" {
		me.DefaultService = me.FixupService(stackRole, me.Default)
	}
	me.Default = ""
	me.ServiceOptions = make(Services, len(me.ServiceIds))
	for i, o := range me.ServiceIds {
		me.ServiceOptions[i] = me.FixupService(stackRole, o)
	}
	me.ServiceIds = nil
}

func (me *RoleService) FixupService(stackRole *StackRole, serviceId ServiceId) (service *Service) {
	service = &Service{
		StackRole: stackRole,
	}
	if me.DefaultService == nil {
		me.DefaultService = &Service{
			StackRole: stackRole,
			Identity: &Identity{
				OrgName: me.OrgName,
			},
		}
	}
	service.Assign(serviceId, me.DefaultService)
	return service
}

func (me *Gears) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *Gears) Refresh() (status stat.Status) {
	var b []byte
	if !me.NeedsRefresh() {
		return status
	}
	for range only.Once {
		cacheDir := me.Gearbox.GetHostConnector().GetCacheDir()
		store := cache.NewCache(cacheDir)

		store.Disable = me.Gearbox.NoCache()
		var ok bool
		b, ok, status = store.Get(GearsKey)
		if ok {
			break
		}
		var sc int
		b, sc, status = util.HttpGet(GearsJsonUrl)
		if status.IsError() || sc != http.StatusOK { // @TODO Bundle these as Assets so we will always have some options
			log.Fatal("Could not download 'gears.json' and no options have previously been stored.")
		}
		status = store.Set(GearsKey, b, "15m")
		if status.IsError() {
			log.Printf(status.Message)
			break
		}
	}
	for range only.Once {
		if status.IsError() {
			break
		}
		status = me.Unmarshal(b)
		if status.IsError() {
			break
		}
		for rs, sr := range me.RoleMap {
			status = sr.Fixup(rs)
			if status.IsError() {
				break
			}
		}
		if status.IsError() {
			break
		}
		for rs, ro := range me.RoleServicesMap {
			sr, ok := me.RoleMap[rs]
			if !ok {
				continue // @TODO Log error here and communicate back to home base
			}
			ro.Fixup(sr)
		}
		me.refreshed = true
	}
	return status
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

func (me *Gears) Unmarshal(b []byte) (status stat.Status) {
	err := json.Unmarshal(b, &me)
	if err != nil {
		// @TODO Provide a link to upgrade once we have that established
		status = stat.NewFailStatus(&stat.Args{
			Message: "failed to unmarshal json from 'gears.json'",
			Help: fmt.Sprintf("Your Gearbox is probably not compatible with the current JSON schema for 'options' at %s. Your Gearbox may need to be updated.",
				GearsJsonUrl,
			),
			Error: err,
		})
	}
	return status
}
