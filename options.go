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

const OptionsKey = "options"
const OptionsJsonUrl = "https://raw.githubusercontent.com/gearboxworks/gearbox/master/assets/options.json"

type Options struct {
	Gearbox *Gearbox       `json:"-"`
	Stacks  StackOptionMap `json:"stacks"`
}

func NewOptions(gb *Gearbox) *Options {
	o := Options{
		Gearbox: gb,
	}
	return &o
}

type StackOptionMap map[StackName]*StackOption
type StackOption struct {
	Authority string             `json:"authority"`
	Roles     StackRoleOptionMap `json:"roles"`
}

type StackRoleOptionMap map[RoleName]*StackRoleOption
type StackRoleOption struct {
	Required bool       `json:"required"`
	Minimum  int        `json:"min"`
	Maximum  int        `json:"max"`
	Default  ServiceId  `json:"default"`
	Options  ServiceIds `json:"options"`
}

func (me *Options) Refresh() (err error) {
	var b []byte
	for range only.Once {
		store := cache.NewCache(me.Gearbox.HostConnector.GetCacheDir())

		store.Disable = me.Gearbox.NoCache()

		b, err = store.Get(OptionsKey)
		if err == nil {
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
