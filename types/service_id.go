package types

import "sort"

type ServiceIds []ServiceId
type ServiceId string

func (me ServiceIds) Sort() {
	sns := make([]string, len(me))
	for i, sn := range me {
		sns[i] = string(sn)
	}
	sort.Strings(sns)
	for i, sn := range sns {
		me[i] = ServiceId(sn)
	}
}
