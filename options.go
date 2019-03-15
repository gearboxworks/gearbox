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

const OptionsJsonUrl = RepoRawBaseUrl + "/master/assets/options.json"
const OptionsKey = "options"

type Options struct {
	Gearbox     *Gearbox    `json:"-"`
	Authorities Authorities `json:"authorities"`
	StackNames  StackNames  `json:"stacks"`
	RoleMap     RoleMap     `json:"roles"`
	OptionsMap  OptionsMap  `json:"role_options"`
	refreshed   bool
}

func NewOptions(gb *Gearbox) *Options {
	o := Options{
		Gearbox: gb,
	}
	return &o
}

type StackOptionMap map[StackName]*StackOption
type StackOption struct {
	Roles OptionsMap `json:"roles"`
}

type ShareableChoices string

const (
	NotShareable    ShareableChoices = "no"
	InStackSharable ShareableChoices = "instack"
	YesShareable    ShareableChoices = "yes"
)

type OptionsMap map[RoleSpec]*RoleOption

func (me OptionsMap) GetStackOptionsMap(stackName StackName) (om OptionsMap, status stat.Status) {
	for range only.Once {
		om = make(OptionsMap, 0)
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

type RoleOption struct {
	*StackRole
	OrgName        OrgName          `json:"org,omitempty"`
	Default        ServiceId        `json:"default,omitempty"`
	Shareable      ShareableChoices `json:"shareable,omitempty"`
	Options        ServiceIds       `json:"options,omitempty"`
	DefaultService *Service         `json:"-"`
	ServiceOptions Services         `json:"-"`
}

func (me *Options) GetStackRoleMap(stackName StackName) RoleMap {
	srs := make(RoleMap, 0)
	for i, r := range me.RoleMap {
		if r.GetStackName() != stackName {
			continue
		}
		srs[i] = r
	}
	return srs
}

func (me *RoleOption) Fixup(stackRole *StackRole) {
	if me.Default != "" {
		me.DefaultService = me.FixupService(stackRole, me.Default)
	}
	me.Default = ""
	me.ServiceOptions = make(Services, len(me.Options))
	for i, o := range me.Options {
		me.ServiceOptions[i] = me.FixupService(stackRole, o)
	}
	me.Options = nil
}

func (me *RoleOption) FixupService(stackRole *StackRole, serviceId ServiceId) (service *Service) {
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

func (me *Options) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *Options) Refresh() (status stat.Status) {
	var b []byte
	var err error
	if !me.NeedsRefresh() {
		return status
	}
	for range only.Once {
		cacheDir := me.Gearbox.HostConnector.GetCacheDir()
		store := cache.NewCache(cacheDir)

		store.Disable = me.Gearbox.NoCache()
		var ok bool
		b, ok, status = store.Get(OptionsKey)
		if ok {
			break
		}
		var sc int
		b, sc, err = util.HttpGet(OptionsJsonUrl)
		if err != nil || sc != http.StatusOK { // @TODO Bundle these as Assets so we will always have some options
			log.Fatal("Could not download 'options.json' and no options have previously been stored.")
		}
		status = store.Set(OptionsKey, b, "15m")
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
			sr.Fixup(rs)
		}
		for rs, ro := range me.OptionsMap {
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

func (me *Options) String() string {
	return string(me.Bytes())
}

func (me *Options) Bytes() []byte {
	bytes, err := json.Marshal(me)
	if err != nil {
		log.Fatal("Could not unserialize 'options.json' cache file.")
	}
	return bytes
}

func (me *Options) Unmarshal(b []byte) (status stat.Status) {
	err := json.Unmarshal(b, &me)
	if err != nil {
		// @TODO Provide a link to upgrade once we have that established
		status = stat.NewFailedStatus(&stat.Args{
			Message: "failed to unmarshal json from 'options.json'",
			Help: fmt.Sprintf("Your Gearbox is probably not compatible with the current JSON schema for 'options' at %s. Your Gearbox may need to be updated.",
				OptionsJsonUrl,
			),
			Error: err,
		})
	}
	return status
}
