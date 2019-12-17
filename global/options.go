package global

import "fmt"

type Options struct {
	IsDebug bool
	NoCache bool
}

func (me *Options) Debug() string {
	return fmt.Sprintf("NoCache:%t,IsDebug:%t", me.IsDebug, me.NoCache)
}
