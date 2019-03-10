package gearbox

import (
	"fmt"
	"regexp"
	"strings"
)

type Stacks []*Stack

type StackMap map[StackName]*Stack

type StackBag map[StackName]interface{}

type StackName string

type Stack struct {
	Name       StackName  `json:"name"`
	Label      string     `json:"label"`
	ServiceMap ServiceMap `json:"services,omitempty"`
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

func (me *Stack) GetServiceMap() ServiceMap {
	mm := make(ServiceMap, len(me.ServiceMap))
	ren := regexp.QuoteMeta(string(me.Label))
	for mt, st := range me.ServiceMap {
		sl := regexp.MustCompile(fmt.Sprintf("^%s", ren)).ReplaceAllString(st.Label, "")
		mm[mt] = &Service{
			Name:        ServiceName(fmt.Sprintf("%s/%s", me.Name, mt)),
			StackName:   me.Name,
			ServiceType: string(mt),
			Label:       st.Label,
			ShortLabel:  strings.Trim(sl, " "),
			Examples:    st.Examples,
			Optional:    st.Optional,
		}
	}
	return mm
}

func GetStackMap() StackMap {
	return StackMap{
		"wordpress": &Stack{
			Name:  "wordpress",
			Label: "WordPress",
			ServiceMap: ServiceMap{
				"webserver": &Service{
					Label:    "WordPress Web Server",
					Examples: []string{"Apache", "Nginx", "Caddy", "Lighttpd"},
				},
				"processvm": &Service{
					Label:    "WordPress Process VM",
					Examples: []string{"PHP", "HHVM"},
				},
				"dbserver": &Service{
					Label:    "WordPress Database Server",
					Examples: []string{"MySQL", "MariaDB", "Percona"},
				},
				"cacheserver": &Service{
					Label:    "WordPress Cache Server",
					Examples: []string{"Redis", "Memcached"},
					Optional: true,
				},
			},
		},
	}
}
