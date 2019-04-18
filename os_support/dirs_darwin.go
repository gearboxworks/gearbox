// +build darwin

package oss

import (
	"fmt"
	"gearbox/types"
)

const SuggestedBasedir types.AbsoluteDir = "Sites"

var NilOsSupport = (*OsSupport)(nil)

var _ OsSupporter = NilOsSupport

type OsSupport struct {
	Nix
}

func (me *OsSupport) GetSuggestedBasedir() types.AbsoluteDir {
	return types.AbsoluteDir(fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		SuggestedBasedir,
	))
}
