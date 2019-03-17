package gearbox

import "fmt"

type GlobalOptions struct {
	IsDebug bool
	NoCache bool
}

//====[GlobalOptions methods]========
func (me *GlobalOptions) Debug() string {
	return fmt.Sprintf("NoCache:%t,IsDebug:%t", me.IsDebug, me.NoCache)
}

//====[Gearbox methods]==============
func (me GearboxObj) IsDebug() bool {
	return me.GlobalOptions.IsDebug
}

func (me GearboxObj) NoCache() bool {
	return me.GlobalOptions.NoCache
}
