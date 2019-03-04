package gearbox

import (
	"fmt"
	"gearbox/util"
	"strings"
)

type Candidates []*Candidate

type Candidate struct {
	BaseDir  string   `json:"base_dir"`
	Path     string   `json:"path"`
	FullPath string   `json:"full_path"`
	Config   *Config  `json:"-"`
	Gearbox  *Gearbox `json:"-"`
}

type CandidateArgs Candidate

func NewCandidate(args *CandidateArgs) *Candidate {
	c := Candidate(*args)
	c.FullPath = c.GetFullPath()
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

func (me *Candidate) GetFullPath() (fp string) {
	fp, _ = ExpandBaseDirPath(me.Gearbox, me.BaseDir, me.Path)
	return fp
}
