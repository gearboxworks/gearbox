package gearbox

import (
	"fmt"
	"regexp"
)

type Stacks []*Stack

type StackMap map[StackName]*Stack

type StackName string

type Stack struct {
	Name    StackName      `json:"name"`
	Label   string         `json:"label"`
	Members StackMemberMap `json:"members,omitempty"`
}

func (me *Stack) String() string {
	return string(me.Name)
}

func (me *Stack) CloneSansMembers() *Stack {
	s := Stack{}
	s.Name = me.Name
	s.Label = me.Label
	return &s
}

func (me *Stack) GetMembers() StackMemberMap {
	mm := make(StackMemberMap, len(me.Members))
	ren := regexp.QuoteMeta(string(me.Name))
	for mt, st := range me.Members {
		sl := regexp.MustCompile(fmt.Sprintf("^%s", ren)).ReplaceAllString(st.Label, "")
		mm[mt] = &StackMember{
			Name:       StackMemberName(fmt.Sprintf("%s/%s", me.Name, mt)),
			StackName:  me.Name,
			MemberType: string(mt),
			Label:      st.Label,
			ShortLabel: sl,
			Examples:   st.Examples,
		}
	}
	return mm
}

func GetStackMap() StackMap {
	return StackMap{
		"wordpress": &Stack{
			Name:  "wordpress",
			Label: "WordPress",
			Members: StackMemberMap{
				"webserver": &StackMember{
					Label:    "WordPress Web Server",
					Examples: []string{"Apache", "Nginx", "Caddy", "Lighttpd"},
				},
				"processvm": &StackMember{
					Label:    "WordPress Process VM",
					Examples: []string{"PHP", "HHVM"},
				},
				"dbserver": &StackMember{
					Label:    "WordPress Database Server",
					Examples: []string{"MySQL", "MariaDB", "Percona"},
				},
				"cacheserver": &StackMember{
					Label:    "WordPress Cache Server",
					Examples: []string{"Redis", "Memcached"},
				},
			},
		},
	}
}
