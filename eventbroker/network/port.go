package network

import "strconv"

type Port string

func (me *Port) String() string {
	return string(*me)
}

func (me *Port) ToInt() int {

	p, _ := strconv.Atoi(me.String())

	return p
}
