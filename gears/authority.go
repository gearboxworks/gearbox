package gears

import "gearbox/types"

type Authorities []*Authority

type Authority struct {
	Domain types.AuthorityDomain `json:"domain"`
}
