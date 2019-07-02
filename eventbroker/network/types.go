package network

import (
	"github.com/grandcat/zeroconf"
)

type Entry zeroconf.ServiceEntry

type (
	Name   = string
	Type   = string
	Domain = string
)

type Text []string

func (me *Text) String() []string {
	return []string(*me)
}
