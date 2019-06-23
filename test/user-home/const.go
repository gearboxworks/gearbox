package testconst

import (
	"gearbox/types"
	"path/filepath"
	"runtime"
)

var UserHomeDir types.Dir

func init() {
	_, fn, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to access filename from runtime.Caller()")
	}
	UserHomeDir = types.Dir(filepath.Dir(fn))
}
