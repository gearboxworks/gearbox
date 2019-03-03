package gearbox

import (
	"fmt"
	"gearbox/util"
	"strings"
)

type Candidates []*Candidate

type Candidate struct {
	BaseDir string  `json:"nickname"`
	Path    string  `json:"path"`
	Config  *Config `json:"-"`
}
type CandidateArgs Candidate

func NewCandidate(args *CandidateArgs) *Candidate {
	c := Candidate(*args)
	return &c
}

func (me *Candidate) GetPotentialHostname() string {
	hostname := me.Path
	if !strings.Contains(hostname, ".") {
		hostname = fmt.Sprintf("%s.local", hostname)
	}
	return strings.ToLower(hostname)
}

func (me *Candidate) GetHostBaseDir() string {
	return me.Config.GetHostBaseDir(me.BaseDir)
}

func (me *Candidate) IsProject() bool {
	return util.FileExists(
		fmt.Sprintf("%s/%s/%s", me.GetHostBaseDir(), me.Path, ProjectFile),
	)
}
