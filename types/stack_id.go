package types

import "sort"

//
// StackId is "authority/stackname"
//
type StackIds []StackId
type StackId string

func (me StackIds) Sort() {
	sns := make([]string, len(me))
	for i, sn := range me {
		sns[i] = string(sn)
	}
	sort.Strings(sns)
	for i, sn := range sns {
		me[i] = StackId(sn)
	}
}
