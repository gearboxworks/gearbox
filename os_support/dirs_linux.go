package oss

// build linux

import (
	"fmt"
)

const SuggestedBasedir = "projects"

var NilOsSupport = (*OsSupport)(nil)

var _ OsSupporter = NilOsSupport

type OsSupport struct {
	Nix
}

func (me *OsSupport) GetSuggestedBasedir() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		SuggestedBasedir,
	)
}
