package gearspec

import (
	"fmt"
	"gearbox/global"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"regexp"
	"strconv"
	"strings"
)

const minStacknameLen = 2
const minRoleLen = 3

type Gearspecs []*Gearspec

func (me Gearspecs) FindById(gsid Identifier) (gs *Gearspec) {
	for range only.Once {
		for _, _gs := range me {
			if _gs.Identifier != gsid {
				continue
			}
			gs = _gs
		}
	}
	return gs
}

type Gearspec struct {
	Identifier      Identifier
	IsRemote        bool
	AuthorityDomain types.AuthorityDomain `json:"authority,omitempty"`
	Stackname       types.Stackname       `json:"stackname,omitempty"`
	Role            types.StackRole       `json:"role,omitempty"`
	Revision        types.Revision        `json:"revision,omitempty"`
}

type Args Gearspec

func NewGearspec() *Gearspec {
	return &Gearspec{}
}

var regexes = struct {
	authority *regexp.Regexp
	httpurl   *regexp.Regexp
}{
	regexp.MustCompile("[^A-Za-z0-9-/.]"),
	regexp.MustCompile("#https?://#"),
}

func (me *Gearspec) ParseString(gearspecid string) (sts status.Status) {
	return me.Parse(Identifier(gearspecid))
}

func (me *Gearspec) ParseStackId(stackid types.StackId) (sts status.Status) {
	gearspecid := fmt.Sprintf("%s/%s", stackid, "dummy")
	sts = me.Parse(Identifier(gearspecid))
	if is.Success(sts) {
		me.Role = ""
	}
	return sts
}

func (me *Gearspec) Parse(gsi Identifier) (sts status.Status) {
	for range only.Once {
		if me == nil {
			panic("gsi.Parse() called when 'gsi' is nil.")
		}
		if regexes.httpurl.Match([]byte(gsi)) {
			sts = me.ParseRemoteGearspec(gsi)
			break
		}
		sts = me.ParseLocalGearspec(gsi)
	}
	return sts
}

func (me *Gearspec) ParseRemoteGearspec(gsi Identifier) (sts status.Status) {
	panic("Remote Gearspecs not yet implemented")
	me.IsRemote = true
	return sts
}

func (me *Gearspec) ParseLocalGearspec(gsi Identifier) (sts status.Status) {
	var err error
	tmp := Gearspec{Identifier: gsi}
	for range only.Once {
		if me == nil {
			panic("gsi.Parse() called when 'gsi' is nil.")
		}
		parts := strings.Split(string(gsi), ":")
		if len(parts) > 1 {
			_, err = strconv.Atoi(parts[1])
			if err != nil {
				sts = status.Wrap(err).
					SetMessage("invalid version in '%s'", gsi).
					SetAllHelp("version must be integer after a colon (':') at end of '%s'", gsi)
				break
			}
			tmp.Revision = types.Revision(parts[1])
		}
		parts = strings.Split(parts[0], "/")
		if len(parts) == 3 && parts[0] == global.DefaultAuthorityDomain {
			parts = []string{parts[1], parts[2]}
		}
		if len(parts) != 2 {
			sts = status.Fail().
				SetMessage("invalid gearspec ID '%s'", gsi).
				SetAllHelp("gearspec ID must contain exactly two (2) slash-separated segments, or be a valid URL, i.e. {stack}/{role}")
			break
		}
		tmp.AuthorityDomain = global.DefaultAuthorityDomain
		tmp.Identifier = gsi
		tmp.Stackname = parts[0]
		if len(tmp.Stackname) < minStacknameLen {
			sts = status.Wrap(err).
				SetMessage("invalid stackname in '%s'", gsi).
				SetAllHelp("stackname must be at least %d characters long", minStacknameLen)
			break
		}
		tmp.Role = parts[1]
		if len(tmp.Role) < minRoleLen {
			sts = status.Wrap(err).
				SetMessage("invalid role in '%s'", gsi).
				SetAllHelp("role must be at least %d characters long", minRoleLen)
		}
	}
	if is.Success(sts) {
		*me = tmp
		me.IsRemote = false
	}
	return sts
}

func (me *Gearspec) GetIdentifier() (id Identifier) {
	var _id string
	for range only.Once {
		if me.IsRemote {
			_id = string(me.Identifier)
			break
		}
		if me.Revision == "" {
			_id = strings.ToLower(fmt.Sprintf("%s/%s", me.Stackname, me.Role))
			break
		}
		_id = strings.ToLower(fmt.Sprintf("%s/%s:%s", me.Stackname, me.Role, me.Revision))
	}
	return Identifier(_id)
}

func (me *Gearspec) String() string {
	return string(me.GetIdentifier())
}

func (me *Gearspec) GetRaw() Identifier {
	return me.Identifier
}

func (me *Gearspec) GetAuthority() types.AuthorityDomain {
	return me.AuthorityDomain
}

func (me *Gearspec) GetStackname() types.Stackname {
	return me.Stackname
}

func (me *Gearspec) GetRole() types.StackRole {
	return me.Role
}

func (me *Gearspec) GetRevision() types.Revision {
	return me.Revision
}

func (me *Gearspec) GetStackId() (sid types.StackId) {
	if me.AuthorityDomain == "" {
		me.AuthorityDomain = global.DefaultAuthorityDomain
	}
	return types.StackId(fmt.Sprintf("%s/%s", me.AuthorityDomain, me.Stackname))
}

func (me *Gearspec) SetStackId(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		tmp := Gearspec{Identifier: Identifier(stackid)}
		if me == nil {
			panic("gearspec.SetStackId() called when stackid is nil.")
		}

		parts := strings.Split(string(stackid), "/")
		if len(parts) < 2 {
			parts = []string{
				global.DefaultAuthorityDomain,
				string(stackid),
			}
		} else if len(parts) > 2 {
			sts = status.Fail().
				SetMessage("invalid stack ID '%s'", stackid).
				SetAllHelp("ID can only have one slash ('/') and it should separate authority from stackname")
			break
		}
		switch len(parts) {
		case 1:
			tmp.AuthorityDomain = global.DefaultAuthorityDomain
			tmp.Stackname = types.Stackname(stackid)
		default:
			tmp.AuthorityDomain = types.AuthorityDomain(parts[0])
			tmp.Stackname = types.Stackname(parts[1])
		}
		*me = tmp
		sts = status.Success("named stack id '%s' set", stackid)
	}
	return sts
}
