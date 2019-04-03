package gearbox

import (
	"gearbox/gears"
	"gearbox/gearspecid"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

type RoleMap map[gsid.Identifier]*gears.StackRole

func (me RoleMap) GetStackRoleMap(stackid types.StackId) (rm RoleMap, sts status.Status) {
	for range only.Once {
		rm = make(RoleMap, 0)
		for rs, r := range me {
			if !strings.HasPrefix(string(rs), string(stackid)) {
				continue
			}
			rm[rs] = r
		}
	}
	return rm, sts
}
