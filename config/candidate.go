package config

import (
	"fmt"
	"gearbox/jsonfile"
	"gearbox/only"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status/is"
	"path/filepath"
	"strings"
)

type Candidates []*Candidate

type Candidate struct {
	Basedir  types.Nickname     `json:"basedir"`
	Path     types.RelativePath `json:"path"`
	FullPath types.AbsoluteDir  `json:"full_path"`
	Config   *Config            `json:"-"`
}

type CandidateArgs Candidate

func NewCandidate(args *CandidateArgs) *Candidate {
	c := Candidate(*args)
	return &c
}

//
// "enmarshal" means "prepare for marshalling
//
func (me *Candidate) Enmarshal() (sts Status) {
	me.FullPath, sts = me.GetFullPath()
	return sts
}

func (me *Candidate) GetPotentialHostname() types.Hostname {
	hostname := types.Hostname(me.Path)
	if !strings.Contains(string(hostname), ".") {
		hostname = types.Hostname(fmt.Sprintf("%s.local", string(hostname)))
	}
	return types.Hostname(strings.ToLower(string(hostname)))
}

func (me *Candidate) GetBasedir() (types.AbsoluteDir, Status) {
	return me.Config.GetBasedir(me.Basedir)
}

func (me *Candidate) IsProject() (ok bool) {
	for range only.Once {
		bd, sts := me.GetBasedir()
		if is.Error(sts) {
			break
		}
		jsfp := jsonfile.GetFilepath(bd, me.Path)
		ok = is.Success(sts) && util.FileExists(jsfp)
	}
	return ok
}

func (me *Candidate) GetFullPath() (fp types.AbsoluteDir, sts Status) {
	for range only.Once {
		fp, sts = me.Config.ExpandBasedirPath(me.Basedir, me.Path)
		if is.Error(sts) {
			break
		}
		fp = types.AbsoluteDir(filepath.FromSlash(string(fp)))
	}
	return fp, sts
}
