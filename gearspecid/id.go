package gsid

import (
	"fmt"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"regexp"
	"strconv"
	"strings"
)

type Id struct {
	raw       Identifier
	Authority types.AuthorityDomain `json:"authority,omitempty"`
	Stackname types.Stackname       `json:"stack,omitempty"`
	Role      types.StackRole       `json:"role,omitempty"`
	Revision  types.Revision        `json:"revision,omitempty"`
}

type Args Id

func NewGearspecId() *Id {
	return &Id{}
}

type reMap map[string]*regexp.Regexp

var re reMap

func init() {
	re = make(reMap, 2)
	re["authority"] = regexp.MustCompile("[^A-Za-z0-9-/.]")
	re["ns_or_r"] = regexp.MustCompile("[^A-Za-z0-9-/]")
}

func (me *Id) ParseString(gearspecid string) (sts status.Status) {
	return me.Parse(Identifier(gearspecid))
}

func (me *Id) ParseStackId(stackid types.StackId) (sts status.Status) {
	gearspecid := fmt.Sprintf("%s/%s", stackid, "dummy")
	sts = me.Parse(Identifier(gearspecid))
	if is.Success(sts) {
		me.Role = ""
	}
	return sts
}

func (me *Id) Parse(gsi Identifier) (sts status.Status) {
	var err error
	tmp := Id{raw: gsi}
	for range only.Once {
		if me == nil {
			panic("gsi.Parse() called when 'gsi' is nil.")
		}
		*me = Id{}
		parts := strings.Split(string(gsi), ":")
		if len(parts) > 1 {
			_, err = strconv.Atoi(parts[1])
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("invalid version in '%s'", gsi),
					Help: fmt.Sprintf("version must be integer after a colon (':') at end of '%s'",
						gsi,
					),
				})
				break
			}
			tmp.Revision = types.Revision(parts[1])
		}
		parts = strings.Split(parts[0], "/")
		if len(parts) == 1 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid gearspec ID '%s'", gsi),
				Help:    "gearspec ID must contain at least two (2) slash-separated segments, i.e. {stack}/{role}",
			})
			break
		}
		tmp.Role = types.StackRole(parts[len(parts)-1])
		sharedHelp := " can only contain letters [a-z], numbers [0-9], dashes ('-')%s or slashes ('/')"
		if strings.Contains(parts[0], ".") {
			tmp.Authority = types.AuthorityDomain(parts[0])
			if re["authority"].MatchString(string(tmp.Authority)) {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("invalid authority '%s' in '%s'", tmp.Authority, gsi),
					Help: fmt.Sprintf("authority '%s' in '%s'%s",
						tmp.Authority,
						gsi,
						fmt.Sprintf(sharedHelp, ", dots ('.')"),
					),
				})
				break
			}
			if len(parts) >= 2 {
				tmp.Stackname = types.Stackname(strings.Join(parts[1:len(parts)-1], "/"))
			}
		} else if len(parts) > 1 {
			tmp.Stackname = types.Stackname(strings.Join(parts[:len(parts)-1], "/"))
		}
		if re["ns_or_r"].MatchString(string(tmp.Stackname)) {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid stack name '%s' in stack ID '%s'",
					tmp.Stackname,
					gsi,
				),
				Help: fmt.Sprintf("stack name '%s'%s",
					tmp.Stackname,
					fmt.Sprintf(sharedHelp, ""),
				),
			})
			break
		}
		if re["ns_or_r"].MatchString(string(tmp.Role)) {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid role '%s' in '%s'",
					tmp.Role,
					gsi,
				),
				Help: fmt.Sprintf("role '%s'%s",
					tmp.Role,
					fmt.Sprintf(sharedHelp, ""),
				),
			})
			break
		}
		if tmp.Authority == "" {
			tmp.Authority = global.DefaultAuthority
		}
	}
	if err == nil {
		*me = tmp
	}
	return sts
}

func (me *Id) String() string {
	var s string
	if me.Authority == "" && me.Stackname == "" && me.Revision == "" {
		s = string(me.Role)
	} else if me.Authority == "" && me.Revision == "" {
		s = fmt.Sprintf("%s/%s", me.Stackname, me.Role)
	} else if me.Authority == "" && me.Stackname == "" {
		s = fmt.Sprintf("%s:%s", me.Role, me.Revision)
	} else if me.Authority == "" {
		s = fmt.Sprintf("%s/%s:%s", me.Stackname, me.Role, me.Revision)
	} else if me.Revision == "" {
		s = fmt.Sprintf("%s/%s/%s", me.Authority, me.Stackname, me.Role)
	} else {
		s = fmt.Sprintf("%s/%s/%s:%s", me.Authority, me.Stackname, me.Role, me.Revision)
	}
	return s
}

func (me *Id) GetRaw() Identifier {
	return me.raw
}

func (me *Id) GetAuthority() types.AuthorityDomain {
	return me.Authority
}

func (me *Id) GetStackname() types.Stackname {
	return me.Stackname
}

func (me *Id) GetRole() types.StackRole {
	return me.Role
}

func (me *Id) GetRevision() types.Revision {
	return me.Revision
}

func (me *Id) GetStackId() types.StackId {
	if me.Authority == "" {
		me.Authority = global.DefaultAuthority
	}
	return types.StackId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}

func (me *Id) SetStackId(stackid types.StackId) (sts status.Status) {
	for range only.Once {
		tmp := Id{raw: Identifier(stackid)}
		if me == nil {
			panic("gearspecid.SetStackId() called when 'spec' is nil.")
		}
		parts := strings.Split(string(stackid), "/")
		if len(parts) < 2 {
			parts = []string{
				global.DefaultAuthority,
				string(stackid),
			}
		} else if len(parts) > 2 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("invalid stack ID '%s'", stackid),
				Help:    "ID can only have one slash ('/') and it should separate authority from stackname",
			})
			break
		}
		switch len(parts) {
		case 1:
			tmp.Stackname = types.Stackname(stackid)
		default:
			tmp.Authority = types.AuthorityDomain(parts[0])
			tmp.Stackname = types.Stackname(parts[1])
		}
		if tmp.Authority == "" {
			tmp.Authority = global.DefaultAuthority
		}
		*me = tmp
		sts = status.Success("named stack id '%s' set", me.GetStackId())
	}
	return sts
}
