package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"strings"
)

type Identity struct {
	raw          string
	group        string
	_type        string
	name         string
	version      *DottedVersion
}
func NewIdentity() (id *Identity) {
	return &Identity{}
}

func (me *Identity) Parse(identity string) (err error) {
	const sharedHelp = "identities can take the form of either " +
		"<group>/<type>/<name>:<version> or just " +
		"<group>/<name>:<version>. Examples might include " +
		"'google/flutter:1.3.8' or 'wordpress/plugins/akismet:4.1.1'"

	var parts []string
	var g string
	var t string
	var n string
	for range only.Once {
		v := NewDottedVersion()
		err = v.Parse(util.After(identity, ":"))
		if err != nil {
			break
		}
		parts = strings.Split(util.Before(identity, ":"), "/")
		switch len(parts) {
		case 1:
			n = parts[0]
		case 2:
			g = parts[0]
			n = parts[1]
		case 3:
			g = parts[0]
			t = parts[1]
			n = parts[2]
		default:
			err = util.AddHelpToError(
				fmt.Errorf("too many slashes ('/') in identity '%s'", identity),
				sharedHelp,
			)
			break
		}
		if n == "" {
			err = util.AddHelpToError(
				fmt.Errorf("name is empty in identity '%s'", identity),
				fmt.Sprintf("%s. So you must have a 'name' such as 'flutter' or 'akismet' in the examples.",
					sharedHelp,
				),
			)
			break
		}
		me.raw = identity
		me.group = g
		me._type = t
		me.name = n
		me.version = v
	}
	return err
}

func (me *Identity) GetId() string {
	return me.String()
}

func (me *Identity) GetRaw() string {
	return me.raw
}

func (me *Identity) GetGroup() string {
	return me.group
}

func (me *Identity) GetType() string {
	return me._type
}

func (me *Identity) GetName() string {
	return me.name
}

func (me *Identity) GetVersion() *DottedVersion {
	if me.version == nil {
		me.version = NewDottedVersion()
	}
	return me.version
}

func (me *Identity) String() (id string) {
	id = me.name
	if me._type != "" {
		id = fmt.Sprintf("%s/%s", me._type, id)
	}
	if me.group != "" {
		id = fmt.Sprintf("%s/%s", me.group, id)
	}
	if me.version != nil && me.version.GetRaw() != "" {
		id = fmt.Sprintf("%s:%s", id, me.version.String())
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
