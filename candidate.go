package gearbox

import (
	"fmt"
	"gearbox/util"
)

type Candidates []*Candidate

type Candidate struct {
	Root *string
	Path string
}

func (me *Candidate) IsProject() bool {
	return util.FileExists(
		fmt.Sprintf("%s/%s/%s", *me.Root, me.Path, ProjectFile),
	)
}
