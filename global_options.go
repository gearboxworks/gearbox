package gearbox

//====[Gearbox methods]==============
func (me *Gearbox) IsDebug() bool {
	return me.GlobalOptions.IsDebug
}

func (me *Gearbox) NoCache() bool {
	return me.GlobalOptions.NoCache
}
