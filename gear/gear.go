package gear

import (
	"fmt"
	"gearbox/global"
	"gearbox/only"
	"gearbox/types"
	"gearbox/util"
	"gearbox/version"
	"github.com/gearboxworks/go-status"
	"strings"
)

type Identifiers []Identifier
type Identifier string

type Gearer interface {
	ParseString(gearid string) (sts status.Status)
	Parse(gearid Identifier) (sts status.Status)
	GetIdentifier() Identifier
	SetRaw(gearid Identifier)
	GetRaw() Identifier
	GetOrgName() types.Orgname
	GetName() types.ProgramName
	GetVersion() *version.Version
	String() (id string)
}

type Gear struct {
	raw     Identifier
	OrgName types.Orgname     `json:"org,omitempty"`
	Program types.ProgramName `json:"program,omitempty"`
	Version *version.Version  `json:"version,omitempty"`
}

func NewGear() (id *Gear) {
	return &Gear{}
}

func (me *Gear) ParseString(gearid string) (sts status.Status) {
	return me.Parse(Identifier(gearid))
}

func (me *Gear) Parse(gearid Identifier) (sts status.Status) {
	const sharedHelp = "identities can take the form of either " +
		"<org>/<type>/<program>:<version> or just " +
		"<org>/<program>:<version>. Examples might include " +
		"'google/flutter:1.3.8' or 'wordpress/plugins/akismet:4.1.1'"

	var parts []string
	var on types.Orgname
	var p types.ProgramName
	var msg string
	var hlp string
	for range only.Once {
		if me == nil {
			panic("gear.Parse() called when 'gear' is nil.")
		}
		v := version.NewVersion()
		sts = v.ParseString(util.After(string(gearid), ":"))
		if status.IsError(sts) {
			break
		}
		before := util.Before(string(gearid), ":")
		if before == "" {
			before = string(gearid)
		}
		parts = strings.Split(before, "/")
		switch len(parts) {
		case 1:
			on = global.DefaultOrgName
			p = types.ProgramName(parts[0])
		case 2:
			on = types.Orgname(parts[0])
			p = types.ProgramName(parts[1])
		default:
			msg = fmt.Sprintf("too many slashes ('/') in gearid '%s'", gearid)
			hlp = sharedHelp
			break
		}
		if p == "" {
			msg = fmt.Sprintf("program is empty in gearid '%s'", gearid)
			hlp = fmt.Sprintf("%s. So you must have a 'name' such as 'flutter' or 'akismet' in the examples.",
				sharedHelp,
			)
			break
		}
		me.raw = gearid
		me.OrgName = on
		me.Program = p
		me.Version = v
	}
	if msg != "" {
		sts = status.Fail(&status.Args{
			Message: msg,
			Help:    hlp,
		})
	}
	return sts
}

func (me *Gear) GetIdentifier() Identifier {
	id := string(me.Program)
	if me.OrgName != "" {
		id = fmt.Sprintf("%s/%s", me.OrgName, id)
	}
	if me.Version != nil && me.Version.GetRaw() != "" {
		id = fmt.Sprintf("%s:%s", id, me.Version.String())
	}
	return Identifier(id)
}

func (me *Gear) String() (id string) {
	return string(me.GetIdentifier())
}

func (me *Gear) SetRaw(gearid Identifier) {
	me.raw = gearid
}

func (me *Gear) GetRaw() Identifier {
	return me.raw
}

func (me *Gear) GetOrgName() types.Orgname {
	return me.OrgName
}

func (me *Gear) GetName() types.ProgramName {
	return me.Program
}

func (me *Gear) GetVersion() *version.Version {
	if me.Version == nil {
		me.Version = version.NewVersion()
	}
	return me.Version
}
