package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"regexp"
	"strconv"
	"strings"
)

type AuthorityDomain string
type Authorities []AuthorityDomain

type Spec struct {
	raw         string
	Authority   AuthorityDomain `json:"authority,omitempty"`
	StackName   StackName       `json:"stack,omitempty"`
	ServiceType string          `json:"type,omitempty"`
	Revision    string          `json:"revision,omitempty"`
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

func (me *Spec) Parse(spec string) (status stat.Status) {
	var err error
	tmp := Spec{raw: spec}
	for range only.Once {
		if me == nil {
			panic("spec.Parse() called when 'spec' is nil.")
		}
		*me = Spec{}
		parts := strings.Split(spec, ":")
		if len(parts) > 1 {
			_, err = strconv.Atoi(parts[1])
			if err != nil {
				status = stat.NewFailedStatus(&stat.Args{
					Error:   err,
					Message: fmt.Sprintf("invalid version in '%s'", spec),
					Help: fmt.Sprintf("version must be integer after a colon (':') at end of '%s'",
						spec,
					),
				})
				break
			}
			tmp.Revision = parts[1]
		}
		sharedHelp := " can only contain letters [a-z], numbers [0-9], dashes ('-')%s or slashes ('/')"
		parts = strings.Split(parts[0], "/")
		tmp.ServiceType = parts[len(parts)-1]
		if strings.Contains(parts[0], ".") {
			tmp.Authority = AuthorityDomain(parts[0])
			if re["host"].MatchString(string(tmp.Authority)) {
				status = stat.NewFailedStatus(&stat.Args{
					Error:   stat.IsStatusError,
					Message: fmt.Sprintf("invalid host '%s' in '%s'", tmp.Authority, spec),
					Help: fmt.Sprintf("host '%s' in '%s'%s",
						tmp.Authority,
						spec,
						fmt.Sprintf(sharedHelp, ", dots ('.')"),
					),
				})
				break
			}
			if len(parts) >= 2 {
				tmp.StackName = StackName(strings.Join(parts[1:len(parts)-1], "/"))
			}
		} else if len(parts) > 1 {
			tmp.StackName = StackName(strings.Join(parts[:len(parts)-1], "/"))
		}
		if re["ns_or_r"].MatchString(string(tmp.StackName)) {
			status = stat.NewFailedStatus(&stat.Args{
				Error: stat.IsStatusError,
				Message: fmt.Sprintf("invalid stack name '%s' in spec '%s'",
					tmp.StackName,
					spec,
				),
				Help: fmt.Sprintf("stack name '%s'%s",
					tmp.StackName,
					fmt.Sprintf(sharedHelp, ""),
				),
			})
			break
		}
		if re["ns_or_r"].MatchString(tmp.ServiceType) {
			status = stat.NewFailedStatus(&stat.Args{
				Error: stat.IsStatusError,
				Message: fmt.Sprintf("invalid role '%s' in '%s'",
					tmp.ServiceType,
					spec,
				),
				Help: fmt.Sprintf("role '%s'%s",
					tmp.ServiceType,
					fmt.Sprintf(sharedHelp, ""),
				),
			})
			break
		}
		if tmp.Authority == "" {
			tmp.Authority = DefaultAuthority
		}
	}
	if err == nil {
		*me = tmp
	}
	return status
}

func (me *Spec) String() string {
	var s string
	if me.Authority == "" && me.StackName == "" && me.Revision == "" {
		s = me.ServiceType
	} else if me.Authority == "" && me.Revision == "" {
		s = fmt.Sprintf("%s/%s", me.StackName, me.ServiceType)
	} else if me.Authority == "" && me.StackName == "" {
		s = fmt.Sprintf("%s:%s", me.ServiceType, me.Revision)
	} else if me.Authority == "" {
		s = fmt.Sprintf("%s/%s:%s", me.StackName, me.ServiceType, me.Revision)
	} else if me.Revision == "" {
		s = fmt.Sprintf("%s/%s/%s", me.Authority, me.StackName, me.ServiceType)
	} else {
		s = fmt.Sprintf("%s/%s/%s:%s", me.Authority, me.StackName, me.ServiceType, me.Revision)
	}
	return s
}

func (me *Spec) GetRaw() string {
	return me.raw
}

func (me *Spec) GetSpec() RoleSpec {
	return RoleSpec(me.String())
}

func (me *Spec) GetAuthority() AuthorityDomain {
	return me.Authority
}

func (me *Spec) GetStackName() StackName {
	return me.StackName
}

func (me *Spec) GetType() string {
	return me.ServiceType
}

func (me *Spec) GetVersion() string {
	return me.Revision
}

func (me *Spec) GetFullStackname() StackName {
	var s StackName
	if me.Authority == "" {
		s = me.StackName
	} else {
		s = StackName(fmt.Sprintf("%s/%s", me.Authority, me.StackName))
	}
	return s
}

func (me *Spec) SetFullStackname(stackName StackName) (status stat.Status) {
	tmp := Spec{raw: string(stackName)}
	if me == nil {
		panic("spec.SetFullStackname() called when 'spec' is nil.")
	}
	parts := strings.Split(string(stackName), "/")
	switch len(parts) {
	case 1:
		tmp.StackName = stackName
	default:
		tmp.Authority = AuthorityDomain(parts[0])
		tmp.StackName = StackName(parts[1])
	}
	if tmp.Authority == "" {
		tmp.Authority = DefaultAuthority
	}
	*me = tmp
	status = stat.NewOkStatus("full stack name '%s' set", me.GetFullStackname())
	return status
}
