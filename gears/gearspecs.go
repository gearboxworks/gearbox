package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"net/http"
	"regexp"
)

const MaxGearsPerSpec = 10

var NilGearspec = (*Gearspec)(nil)
var _ gearspec.Gearspecer = NilGearspec

type GearspecMap map[gearspec.Identifier]*Gearspec
type Gearspecs []*Gearspec

func (me Gearspecs) Find(gsid gearspec.Identifier) (gs *Gearspec, sts Status) {
	for _, _gs := range me {
		if _gs.GearspecId != gsid {
			continue
		}
		gs = _gs
		sts = status.Success("Gearspec '%s' found", gsid)
		break
	}
	if gs == nil {
		sts = status.Fail().
			SetMessage("Gearspec '%s' not found", gsid).
			SetHttpStatus(http.StatusNotFound)
	}
	return gs, sts
}

func (me Gearspecs) GetMap() (gsm GearspecMap) {
	gsm = make(GearspecMap, len(me))
	for _, gs := range me {
		gsm[gs.GearspecId] = gs
	}
	return gsm
}

type Gearspec struct {
	GearspecId      gearspec.Identifier    `json:"gearspec_id"`
	AuthorityDomain types.AuthorityDomain  `json:"authority"`
	Stackname       types.Stackname        `json:"stackname"`
	Specname        types.Specname         `json:"specname"`
	Revision        types.Revision         `json:"revision,omitempty"`
	Label           string                 `json:"label,omitempty"`
	Described       string                 `json:"described,omitempty"`
	Examples        []string               `json:"examples,omitempty"`
	Optional        bool                   `json:"optional,omitempty"`
	Shareable       global.ShareableChoice `json:"shareable"`
	DefaultGearId   service.Identifier     `json:"default,omitempty"`
	Minimum         int                    `json:"min,omitempty"`
	Maximum         int                    `json:"max,omitempty"`
	Gears           Gears                  `json:"-"`

	stackid types.StackId
}

func (me *Gearspec) GetIdentifier() gearspec.Identifier {
	return me.GearspecId
}

var httpRegex = regexp.MustCompile("#https?://#")

func (me *Gearspec) GetIsRemote() bool {
	return httpRegex.Match([]byte(me.GearspecId))
}

func (me *Gearspec) GetAuthorityDomain() types.AuthorityDomain {
	return me.AuthorityDomain
}

func (me *Gearspec) GetStackname() types.Stackname {
	return me.Stackname
}

func (me *Gearspec) GetSpecname() types.Specname {
	return me.Specname
}

func (me *Gearspec) GetRevision() types.Revision {
	return me.Revision
}

type GearspecArgs Gearspec

func NewGearspec() *Gearspec {
	return &Gearspec{}
}

func (me Gearspecs) FilterByNamedStack(stackid types.StackId) (nsrs Gearspecs, sts status.Status) {
	ns := NewNamedStack(stackid)
	stackid = ns.GetIdentifier()
	nsrs = make(Gearspecs, 0)
	for _, gs := range me {
		if gs.GetStackId() != stackid {
			continue
		}
		nsrs = append(nsrs, gs)
	}
	return nsrs, sts
}

func (me *Gearspec) GetDefaultGear() (ds *Gear) {
	for range only.Once {
		if me.DefaultGearId == ZeroString {
			break
		}
		ds = NewGear()
		sts := ds.Parse(me.DefaultGearId)
		if is.Error(sts) {
			sts.Log()
		}
	}
	return ds
}

func (me *Gearspec) GetStackId() types.StackId {
	for range only.Once {
		if me.stackid != ZeroString {
			break
		}
		if me.AuthorityDomain == ZeroString {
			me.AuthorityDomain = global.DefaultAuthorityDomain
		}
		me.stackid = types.StackId(fmt.Sprintf("%s/%s", me.AuthorityDomain, me.Stackname))
	}
	return me.stackid
}

func (me *Gearspec) GetGearspecId() gearspec.Identifier {
	for range only.Once {
		if me.GearspecId != ZeroString {
			break
		}
		gs := gearspec.NewGearspec()
		gs.AuthorityDomain = me.AuthorityDomain
		gs.Stackname = me.Stackname
		gs.Specname = me.Specname
		gs.Revision = me.Revision
		me.GearspecId = gs.GetIdentifier()
	}
	return me.GearspecId
}

func (me *Gearspec) String() string {
	return fmt.Sprintf("%#v", me)
}

func (me *Gearspec) Parse(gsi gearspec.Identifier) (sts status.Status) {
	for range only.Once {
		gs := gearspec.NewGearspec()
		sts = gs.Parse(gsi)
		if is.Error(sts) {
			break
		}
		me.AuthorityDomain = gs.AuthorityDomain
		me.Stackname = gs.Stackname
		me.Specname = gs.Specname
		me.Revision = gs.Revision
		me.stackid = gs.GetStackId()
		me.GearspecId = gs.Identifier
		sts = status.Success("stack role '%s' successfully parsed", gsi)
	}
	return sts
}

func (me *Gearspec) Fixup() (sts status.Status) {
	for range only.Once {
		save := *me
		if me.GearspecId != ZeroString {
			sts = me.Parse(me.GearspecId)
			if status.IsError(sts) {
				break
			}
		}
		if me.Stackname == ZeroString {
			sts = status.Fail().SetMessage("stackname cannot be null for stack role: %s", me)
			break
		}
		if me.Specname == ZeroString {
			sts = status.Fail().SetMessage("role cannot be null for stack role: %s", me)
			break
		}
		if save.Stackname != ZeroString {
			me.Stackname = save.Stackname
		}
		if save.Specname != ZeroString {
			me.Specname = save.Specname
		}
		if me.stackid == ZeroString {
			me.stackid = me.GetStackId()
		}
		if me.GearspecId == ZeroString {
			me.GearspecId = me.GetGearspecId()
		}
		if me.Minimum == 0 && !me.Optional {
			me.Minimum = 1
		}
		if me.Maximum == 0 {
			me.Maximum = MaxGearsPerSpec
		}
	}
	return sts
}
