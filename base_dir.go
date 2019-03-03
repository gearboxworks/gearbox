package gearbox

import (
	"errors"
	"fmt"
	"gearbox/only"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"strings"
)

const DefaultBaseDirNickname = "primary"

type BaseDirMap map[string]*BaseDir

type BaseDirs []*BaseDir

type BaseDir struct {
	Nickname string `json:"nickname"`
	HostDir  string `json:"host_dir"`
	VmDir    string `json:"vm_dir"`
	Error    error  `json:"-"`
}
type BaseDirArgs BaseDir

func NewBaseDir(hostDir string, args ...*BaseDirArgs) *BaseDir {
	var _args *BaseDirArgs
	if len(args) == 0 {
		_args = &BaseDirArgs{HostDir: hostDir}
	} else {
		_args = args[0]
	}
	bd := BaseDir(*_args)
	bd.HostDir = hostDir
	bd.Initialize()
	return &bd
}

func (me *BaseDir) Initialize() {
	for range only.Once {
		if me.HostDir == "" {
			me.Error = errors.New("BaseDir.HostDir has no value")
			break
		}
		if strings.HasPrefix(me.HostDir, "~") {
			var err error
			me.HostDir, err = homedir.Expand(me.HostDir)
			if err != nil {
				me.Error = err
			}
			break
		}
		if me.Nickname == "" {
			me.Nickname = filepath.Base(me.HostDir)
		}
		if me.VmDir == "primary" {
			me.VmDir = vmBaseDir
			break
		}
		if me.VmDir == "" || me.VmDir == vmBaseDir {
			me.VmDir = fmt.Sprintf("%s/%s", vmBaseDir, me.Nickname)
		}
	}
}

// @TODO Delegate responsibility for the VM dir to the VM
//func getVmSubdirFromHostDir(vmRootDir, hostDir string) (vmSubdir string) {
//	base := filepath.Base(hostDir)
//	var index int
//	vmSubdir = fmt.Sprintf("%s/%s", vmRootDir, base )
//	for !util.DirExists(vmSubdir) {
//		index++
//		vmSubdir = fmt.Sprintf("%s/%s%d", vmRootDir, base, index)
//	}
//	return vmSubdir
//}
