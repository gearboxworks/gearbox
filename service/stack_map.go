package service

import (
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
)

type StackMap map[gearspec.Identifier]*ServicerProxy

func (me StackMap) GetNamedStackIds() (nsids types.StackIds, sts status.Status) {
	for range only.Once {
		names := util.NewUniqueStrings(len(me))
		for gs := range me {
			var nsid types.StackId
			nsid, sts = gs.GetNamedStackId()
			if is.Error(sts) {
				break
			}
			names[string(nsid)] = true
		}
		if is.Error(sts) {
			break
		}
		nsids = make(types.StackIds, len(names))
		for _, s := range names.ToSlice() {
			nsids = append(nsids, types.StackId(s))
		}
	}
	return nsids, sts
}
