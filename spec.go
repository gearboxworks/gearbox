package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"regexp"
	"strconv"
	"strings"
)

type Spec struct {
	raw       string
	Host      string `json:"host,omitempty"`
	Namespace string `json:"named_stack,omitempty"`
	Role      string `json:"role"`
	Revision  string `json:"revision,omitempty"`
}

type SpecArgs Spec

func NewSpec() *Spec {
	return &Spec{}
}

type reMap map[string]*regexp.Regexp

var re reMap

func init() {
	re = make(reMap, 2)
	re["host"] = regexp.MustCompile("[^A-Za-z0-9/.]")
	re["ns_or_r"] = regexp.MustCompile("[^A-Za-z0-9/]")
}

func (me *Spec) GetStackName() (name string) {
	return me.Namespace
}

func (me *Spec) Parse(spec string) (err error) {
	tmp := Spec{raw: spec}
	for range only.Once {
		*me = Spec{}
		parts := strings.Split(spec, ":")
		if len(parts) > 1 {
			_, err = strconv.Atoi(parts[1])
			if err != nil {
				err = util.AddHelpToError(
					fmt.Errorf("invalid version in '%s'", spec),
					fmt.Sprintf(
						"version must be integer after a colon (':') at end of '%s'",
						spec,
					),
				)
				break
			}
			tmp.Revision = parts[1]
		}
		sharedHelp := " can only contain letters [a-z], numbers [0-9], dashes ('-')%s or slashes ('/')"
		parts = strings.Split(parts[0], "/")
		tmp.Role = parts[len(parts)-1]
		if strings.Contains(parts[0], ".") {
			tmp.Host = parts[0]
			if re["host"].MatchString(tmp.Host) {
				err = util.AddHelpToError(
					fmt.Errorf("invalid host '%s' in '%s'", tmp.Host, spec),
					fmt.Sprintf("host '%s' in '%s'%s",
						tmp.Host,
						spec,
						fmt.Sprintf(sharedHelp, ", dots ('.')"),
					),
				)
				break
			}
			if len(parts) >= 2 {
				tmp.Namespace = strings.Join(parts[1:len(parts)-1], "/")
			}
		} else if len(parts) > 1 {
			tmp.Namespace = strings.Join(parts[:len(parts)-1], "/")
		}
		if re["ns_or_r"].MatchString(tmp.Namespace) {
			err = util.AddHelpToError(
				fmt.Errorf("invalid namespace '%s' in spec '%s'", tmp.Namespace, spec),
				fmt.Sprintf("namespace '%s'%s", tmp.Namespace, fmt.Sprintf(sharedHelp, "")),
			)
			break
		}
		if re["ns_or_r"].MatchString(tmp.Role) {
			err = util.AddHelpToError(
				fmt.Errorf("invalid role '%s' in '%s'", tmp.Role, spec),
				fmt.Sprintf("role '%s'%s", tmp.Role, fmt.Sprintf(sharedHelp, "")),
			)
			break
		}
	}
	if err == nil {
		*me = tmp
	}
	return err
}

func (me *Spec) String() string {
	var s string
	if me.Host == "" && me.Namespace == "" && me.Revision == "" {
		s = me.Role
	} else if me.Host == "" && me.Revision == "" {
		s = fmt.Sprintf("%s/%s", me.Namespace, me.Role)
	} else if me.Host == "" && me.Namespace == "" {
		s = fmt.Sprintf("%s:%s", me.Role, me.Revision)
	} else if me.Host == "" {
		s = fmt.Sprintf("%s/%s:%s", me.Namespace, me.Role, me.Revision)
	} else if me.Revision == "" {
		s = fmt.Sprintf("%s/%s/%s", me.Host, me.Namespace, me.Role)
	} else {
		s = fmt.Sprintf("%s/%s/%s:%s", me.Host, me.Namespace, me.Role, me.Revision)
	}
	return s
}

func (me *Spec) GetRaw() string {
	return me.raw
}

func (me *Spec) GetSpec() string {
	return me.String()
}

func (me *Spec) GetHost() string {
	return me.Host
}

func (me *Spec) GetNamespace() string {
	return me.Namespace
}

func (me *Spec) GetRole() string {
	return me.Role
}

func (me *Spec) GetVersion() string {
	return me.Revision
}
