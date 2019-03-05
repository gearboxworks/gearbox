package gearbox

type GlobalOptions struct {
	IsDebug bool
}

func (me *Gearbox) IsDebug() bool {
	return me.Options.IsDebug
}
