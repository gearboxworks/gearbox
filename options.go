package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/cache"
	"gearbox/only"
	"gearbox/util"
	"log"
	"net/http"
)

const OptionsJsonUrl = RepoRawBaseUrl + "/master/assets/options.json"
const OptionsKey = "options"

type Options struct {
	Gearbox       *Gearbox      `json:"-"`
	Authorities   Authorities   `json:"authorities"`
	StackNames    StackNames    `json:"stacks"`
	RoleMap       RoleMap       `json:"roles"`
	RoleOptionMap RoleOptionMap `json:"role_options"`
}

func NewOptions(gb *Gearbox) *Options {
	o := Options{
		Gearbox: gb,
	}
	return &o
}

type StackOptionMap map[StackName]*StackOption
type StackOption struct {
	Roles RoleOptionMap `json:"roles"`
}

type RoleOptionMap map[RoleSpec]*RoleOption
type RoleOption struct {
	*StackRole
	OrgName        OrgName    `json:"org,omitempty"`
	Default        ServiceId  `json:"default,omitempty"`
	Options        ServiceIds `json:"options,omitempty"`
	DefaultService *Service   `json:"-"`
	ServiceOptions Services   `json:"-"`
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

func (me *Options) Refresh() (err error) {
	var b []byte
	for range only.Once {
		store := cache.NewCache(me.Gearbox.HostConnector.GetCacheDir())

		store.Disable = me.Gearbox.NoCache()
		var ok bool
		b, ok, err = store.Get(OptionsKey)
		if ok {
			break
		}
		var sc int
		b, sc, err = util.HttpGet(OptionsJsonUrl)
		if err != nil || sc != http.StatusOK {
			// @TODO Bundle these as Assets so we will always have some options
			log.Fatal("Could not download 'options.json' and no options have previously been stored.")
		}
		err = store.Set(OptionsKey, b, "15m")
		if err != nil {
			log.Printf("Could not cache downloaded 'options.json': %s",
				err.Error(),
			)
		}
	}
	err = me.Unmarshal(b)
	if err != nil {
		// @TODO: This needs to become a lot more robust
		//        We should be able to process older versions
		err = util.AddHelpToError(
			err,
			fmt.Sprintf("Your Gearbox is not compatible with the current JSON schema for 'options' at %s. Chances are your Gearbox needs to be updated.",
				OptionsJsonUrl, // @TODO Provide a link to upgrade in text above
			),
		)
	}
	for rs, sr := range me.RoleMap {
		sr.Fixup(rs)
	}
	for rs, ro := range me.RoleOptionMap {
		sr, ok := me.RoleMap[rs]
		if !ok {
			// @TODO Log error here and communicate back to home base
			continue
		}
		ro.Fixup(sr)
	}
	return err
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

func (me *Options) Unmarshal(b []byte) (err error) {
	return json.Unmarshal(b, &me)
}
