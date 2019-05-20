package daemon

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"
	"os"
)

func IsParentInit() (bool) {

	ppid := os.Getppid()
	if ppid == 1 {
		return true
	}

	return false
}


func IsRunning() (bool) {

	fmt.Printf("PPID:%v:\n", IsParentInit())

	foo1, _ := process.Pids()
	for i, p := range foo1 {
		fmt.Printf("process.Pids:%v:	%v:\n", i, p)
	}

	foo2, _ := process.Processes()
	for _, p := range foo2 {
		c, _ := p.Cmdline()
		fmt.Printf("process.Processes:%v:	'%s'\n", p.Pid, c)
	}

	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %v\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %v\n", miscStat.ProcsRunning)

	return false
}
