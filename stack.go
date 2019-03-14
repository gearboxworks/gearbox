package gearbox

import (
	"fmt"
	"regexp"
	"strings"
)

type Stacks []*Stack

type StackMap map[RoleName]*Stack

type StackBag map[RoleName]interface{}

type StackName string

type StackNames []string

type Stack struct {
	Name    StackName `json:"name"`
	Label   string    `json:"label"`
	RoleMap RoleMap   `json:"services,omitempty"`
}

func (me *Stack) String() string {
	return string(me.Name)
}

func (me *Stack) CloneSansServices() *Stack {
	s := Stack{}
	s.Name = me.Name
	s.Label = me.Label
	return &s
}

func (me *Stack) GetRoleMap() RoleMap {
	mm := make(RoleMap, len(me.RoleMap))
	ren := regexp.QuoteMeta(string(me.Label))
	for mt, r := range me.RoleMap {
		sl := regexp.MustCompile(fmt.Sprintf("^%s", ren)).ReplaceAllString(r.Label, "")
		mm[mt] = &StackRole{
			Role:        RoleName(fmt.Sprintf("%s/%s", me.Name, mt)),
			StackName:   me.Name,
			ServiceType: string(mt),
			Label:       r.Label,
			ShortLabel:  strings.Trim(sl, " "),
			Examples:    r.Examples,
			Optional:    r.Optional,
		}
	}
	return mm
}

func GetStackMap() StackMap {
	return StackMap{
		"wordpress": &Stack{
			Name:  "wordpress",
			Label: "WordPress",
			RoleMap: RoleMap{
				"webserver": &StackRole{
					Label:    "WordPress Web Server",
					Examples: []string{"Apache", "Nginx", "Caddy", "Lighttpd"},
				},
				"processvm": &StackRole{
					Label:    "WordPress Process VM",
					Examples: []string{"PHP", "HHVM"},
				},
				"dbserver": &StackRole{
					Label:    "WordPress Database Server",
					Examples: []string{"MySQL", "MariaDB", "Percona"},
				},
				"cacheserver": &StackRole{
					Label:    "WordPress Cache Server",
					Examples: []string{"Redis", "Memcached"},
					Optional: true,
				},
			},
		},
	}
}
