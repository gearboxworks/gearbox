package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"net/http"
)

type Stacks []*Stack

type StackMap map[RoleSpec]*Stack

type ServiceBag map[RoleSpec]interface{}

type StackName string

type StackNames []StackName

type Stack struct {
	Name            StackName       `json:"name"`
	RoleMap         RoleMap         `json:"roles,omitempty"`
	RoleServicesMap RoleServicesMap `json:"role_services,omitempty"`
	Gearbox         *Gearbox        `json:"-"`
	refreshed       bool
}

func NewStack(gb *Gearbox, name StackName) *Stack {
	return &Stack{
		Name:    name,
		Gearbox: gb,
	}
}

func (me *Stack) String() string {
	return string(me.Name)
}

func (me *Stack) LightweightClone() *Stack {
	return NewStack(me.Gearbox, me.Name)
}

func (me *Stack) GetRoleMap() RoleMap {
	mm := make(RoleMap, len(me.RoleMap))
	for mt, r := range me.RoleMap {
		mm[mt] = &StackRole{
			RoleSpec:   RoleSpec(fmt.Sprintf("%s/%s", me.Name, mt)),
			Label:      r.Label,
			ShortLabel: r.ShortLabel,
			Examples:   r.Examples,
			Optional:   r.Optional,
			Spec: &Spec{
				Authority:   DefaultAuthority,
				StackName:   me.Name,
				ServiceType: string(mt),
			},
		}
	}
	return mm
}

func GetStackMap() StackMap {

	return StackMap{
		"wordpress": &Stack{
			Name: "wordpress",
			RoleMap: RoleMap{
				"webserver": &StackRole{
					Label:      "WordPress Web Server",
					ShortLabel: "Web Server",
					Examples:   []string{"Apache", "Nginx", "Caddy", "Lighttpd"},
				},
				"processvm": &StackRole{
					Label:      "WordPress Process VM",
					ShortLabel: "Process VM",
					Examples:   []string{"PHP", "HHVM"},
				},
				"dbserver": &StackRole{
					Label:      "WordPress Database Server",
					ShortLabel: "DB Server",
					Examples:   []string{"MySQL", "MariaDB", "Percona"},
				},
				"cacheserver": &StackRole{
					Label:      "WordPress Cache Server",
					ShortLabel: "Cache Server",
					Examples:   []string{"Redis", "Memcached"},
					Optional:   true,
				},
			},
		},
	}
}

func (me *Stack) GetDefaultServices() (sm ServiceMap, status stat.Status) {
	sm = make(ServiceMap, 0)
	me.Refresh()
	for gs, s := range me.RoleServicesMap {
		if s.DefaultService == nil {
			continue
		}
		sm[gs] = s.DefaultService
	}
	return sm, status
}

func (me *Stack) NeedsRefresh() bool {
	return !me.refreshed
}

func (me *Stack) Refresh() (status stat.Status) {
	for range only.Once {
		if !me.NeedsRefresh() {
			break
		}
		gears := NewGears(me.Gearbox)
		status := gears.Refresh()
		if status.IsError() {
			break
		}
		var srm RoleMap
		srm, status = gears.RoleMap.GetStackRoleMap(me.Name)
		if status.IsError() {
			break
		}
		me.RoleMap = srm

		var sro RoleServicesMap
		sro, status = gears.RoleServicesMap.GetStackServicesMap(me.Name)
		if status.IsError() {
			break
		}
		me.RoleServicesMap = sro
		me.refreshed = true
	}
	if !status.IsError() {
		status = stat.NewOkStatus("named stack '%s' refreshed", me.Name)
	}
	return status
}

func GetFullStackName(stackName StackName) (fsn StackName, status stat.Status) {
	for range only.Once {
		spec := NewSpec()
		status := spec.SetFullStackname(stackName)
		if status.IsError() {
			break
		}
		fsn = spec.GetFullStackname()
	}
	return fsn, status
}

func FindNamedStack(gb *Gearbox, stackName StackName) (stack *Stack, status stat.Status) {
	var tmp *Stack
	for range only.Once {
		status = ValidateStackName(gb, stackName)
		if status.IsError() {
			break
		}
		stackName, status = GetFullStackName(stackName)
		if status.IsError() {
			break
		}
		tmp = NewStack(gb, stackName)
		status = tmp.Refresh()
		if status.IsError() {
			break
		}
	}
	if !status.IsError() {
		stack = &Stack{}
		*stack = *tmp
	}
	return stack, status
}

func ValidateStackName(gb *Gearbox, stackName StackName) (status stat.Status) {
	for range only.Once {
		spec := NewSpec()
		status := spec.SetFullStackname(stackName)
		if status.IsError() {
			break
		}
		gears := NewGears(gb)
		status = gears.Refresh()
		if status.IsError() {
			break
		}
		stackName := spec.GetFullStackname()
		var ok bool
		for _, sn := range gears.StackNames {
			if sn == stackName {
				ok = true
				break
			}
		}
		if !ok {
			status = stat.NewFailedStatus(&stat.Args{
				Message:    fmt.Sprintf("stack '%s' not found", stackName),
				HttpStatus: http.StatusNotFound,
				Help:       fmt.Sprintf("see valid stack names at %s", GearsJsonUrl),
			})
		} else {
			status = stat.NewOkStatus("validated stack name '%s'", stackName)
		}
	}
	return status
}
