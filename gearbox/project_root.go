package gearbox

import (
	"fmt"
	"path/filepath"
)

type ProjectRoots []*ProjectRoot

type ProjectRoot struct {
	HostDir string `json:"host_dir"`
	VmDir   string `json:"vm_dir"`
}

func NewProjectRoot(vmRootDir, hostDir string) *ProjectRoot {
	pr := ProjectRoot{
		HostDir: hostDir,
		VmDir:   getVmSubdirFromHostDir(vmRootDir, hostDir),
	}
	return &pr
}

func getVmSubdirFromHostDir(vmRootDir, hostDir string) (vmSubdir string) {
	base := filepath.Base(hostDir)
	return fmt.Sprintf("%s/%s", vmRootDir, base)
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
