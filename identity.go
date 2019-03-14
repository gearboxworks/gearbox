package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"strings"
)

type OrgName string
type Identity struct {
	raw     string
	OrgName OrgName        `json:"org,omitempty"`
	Type    string         `json:"type,omitempty"`
	Program string         `json:"program,omitempty"`
	Version *DottedVersion `json:"version,omitempty"`
}

func NewIdentity() (id *Identity) {
	return &Identity{}
}

func (me *Identity) Parse(identity string) (err error) {
	const sharedHelp = "identities can take the form of either " +
		"<org>/<type>/<program>:<version> or just " +
		"<org>/<program>:<version>. Examples might include " +
		"'google/flutter:1.3.8' or 'wordpress/plugins/akismet:4.1.1'"

	var parts []string
	var on OrgName
	var t string
	var p string
	for range only.Once {
		if me == nil {
			panic("identity.Parse() called when 'identity' is nil.")
		}
		v := NewDottedVersion()
		err = v.Parse(util.After(identity, ":"))
		if err != nil {
			break
		}
		before := util.Before(identity, ":")
		if before == "" {
			before = identity
		}
		parts = strings.Split(before, "/")
		switch len(parts) {
		case 1:
			p = parts[0]
		case 2:
			on = OrgName(parts[0])
			p = parts[1]
		case 3:
			on = OrgName(parts[0])
			t = parts[1]
			p = parts[2]
		default:
			err = util.AddHelpToError(
				fmt.Errorf("too many slashes ('/') in identity '%s'", identity),
				sharedHelp,
			)
			break
		}
		if p == "" {
			err = util.AddHelpToError(
				fmt.Errorf("program is empty in identity '%s'", identity),
				fmt.Sprintf("%s. So you must have a 'name' such as 'flutter' or 'akismet' in the examples.",
					sharedHelp,
				),
			)
			break
		}
		me.raw = identity
		me.OrgName = on
		me.Type = t
		me.Program = p
		me.Version = v
	}
	return err
}

func (me *Identity) GetId() ServiceId {
	return ServiceId(me.String())
}

func (me *Identity) GetRaw() string {
	return me.raw
}

func (me *Identity) GetOrgName() OrgName {
	return me.OrgName
}

func (me *Identity) GetType() string {
	return me.Type
}

func (me *Identity) GetName() string {
	return me.Program
}

func (me *Identity) GetVersion() *DottedVersion {
	if me.Version == nil {
		me.Version = NewDottedVersion()
	}
	return me.Version
}

func (me *Identity) String() (id string) {
	id = me.Program
	if me.Type != "" {
		id = fmt.Sprintf("%s/%s", me.Type, id)
	}
	if me.OrgName != "" {
		id = fmt.Sprintf("%s/%s", me.OrgName, id)
	}
	if me.Version != nil && me.Version.GetRaw() != "" {
		id = fmt.Sprintf("%s:%s", id, me.Version.String())
	}
	return id
}

//func chkParsed(l *PreviousIdentity) {
//	if ! l.parsed {
//		// See: https://stackoverflow.com/a/25927915/102699
//		// See also: https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
//		// See also: https://lawlessguy.wordpress.com/2016/04/17/display-file-function-and-line-number-in-go-golang/
//		pc := make([]uintptr, 1)
//		runtime.Callers(2, pc)
//		f := runtime.FuncForPC(pc[0])
//		file, line := f.FileLine(pc[0])
//		msg := "Parse() not yet called, in %s:%d"
//		msg = fmt.Sprintf(msg, file, line)
//		panic(errors.New(msg))
//	}
//}
//
//
//
