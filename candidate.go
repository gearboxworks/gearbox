package gearbox

import (
	"fmt"
	"gearbox/util"
)

type Candidates []*Candidate

type Candidate struct {
	Group int     `json:"group"`
	Path  string  `json:"path"`
	Root  *string `json:"-"`
}

func (me *Candidate) IsProject() bool {
	return util.FileExists(
		fmt.Sprintf("%s/%s/%s", *me.Root, me.Path, ProjectFile),
	)
}
