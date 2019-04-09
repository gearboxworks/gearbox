package gearid

import (
	"fmt"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"gearbox/util"
	"gearbox/version"
	"strings"
)

type GearIdentifiers []GearIdentifier
type GearIdentifier string

type GearIder interface {
	ParseString(gearid string) (sts status.Status)
	Parse(gearid GearIdentifier) (sts status.Status)
	GetIdentifier() GearIdentifier
	SetRaw(gearid GearIdentifier)
	GetRaw() GearIdentifier
	GetOrgName() types.OrgName
	GetType() types.ServiceType
	GetName() types.ProgramName
	GetVersion() *version.Version
	String() (id string)
}

type GearId struct {
	raw     GearIdentifier
	OrgName types.OrgName     `json:"org,omitempty"`
	Type    types.ServiceType `json:"type,omitempty"`
	Program types.ProgramName `json:"program,omitempty"`
	Version *version.Version  `json:"version,omitempty"`
}

func NewGearId() (id *GearId) {
	return &GearId{}
}

func (me *GearId) ParseString(gearid string) (sts status.Status) {
	return me.Parse(GearIdentifier(gearid))
}

func (me *GearId) Parse(gearid GearIdentifier) (sts status.Status) {
	const sharedHelp = "identities can take the form of either " +
		"<org>/<type>/<program>:<version> or just " +
		"<org>/<program>:<version>. Examples might include " +
		"'google/flutter:1.3.8' or 'wordpress/plugins/akismet:4.1.1'"

	var parts []string
	var on types.OrgName
	var t types.ServiceType
	var p types.ProgramName
	var msg string
	var hlp string
	for range only.Once {
		if me == nil {
			panic("gearid.Parse() called when 'gearid' is nil.")
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
			on = types.OrgName(parts[0])
			p = types.ProgramName(parts[1])
		case 3:
			on = types.OrgName(parts[0])
			t = types.ServiceType(parts[1])
			p = types.ProgramName(parts[2])
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
		me.Type = t
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

func (me *GearId) GetIdentifier() GearIdentifier {
	id := string(me.Program)
	if me.Type != "" {
		id = fmt.Sprintf("%s/%s", me.Type, id)
	}
	if me.OrgName != "" {
		id = fmt.Sprintf("%s/%s", me.OrgName, id)
	}
	if me.Version != nil && me.Version.GetRaw() != "" {
		id = fmt.Sprintf("%s:%s", id, me.Version.String())
	}
	return GearIdentifier(id)
}

func (me *GearId) String() (id string) {
	return string(me.GetIdentifier())
}

func (me *GearId) SetRaw(gearid GearIdentifier) {
	me.raw = gearid
}

func (me *GearId) GetRaw() GearIdentifier {
	return me.raw
}

func (me *GearId) GetOrgName() types.OrgName {
	return me.OrgName
}

func (me *GearId) GetType() types.ServiceType {
	return me.Type
}

func (me *GearId) GetName() types.ProgramName {
	return me.Program
}

func (me *GearId) GetVersion() *version.Version {
	if me.Version == nil {
		me.Version = version.NewVersion()
	}
	return me.Version
}
