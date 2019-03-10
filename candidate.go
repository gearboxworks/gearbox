package gearbox

import (
	"fmt"
	"gearbox/util"
	"strings"
)

type Candidates []*Candidate

type Candidate struct {
	Basedir  string   `json:"basedir"`
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

func (me *Candidate) GetHostBasedir() (string,error) {
	return me.Config.GetHostBasedir(me.Basedir)
}

func (me *Candidate) IsProject() bool {
	bd,err := me.GetHostBasedir()
	return err == nil && util.FileExists(
		fmt.Sprintf("%s/%s/%s", bd, me.Path, ProjectFilename),
	)
}

func (me *Candidate) GetFullPath() (fp string) {
	fp, _ = ExpandHostBasedirPath(me.Gearbox, me.Basedir, me.Path)
	return fp
}
