package test

import (
	"gearbox/types"
	"strings"
)

func ParseRelativePath(basedir types.AbsoluteDir, fulldir types.AbsoluteFilepath) (path types.RelativePath) {
	if strings.HasPrefix(string(fulldir), string(basedir)) {
		path = types.RelativePath(string([]byte(fulldir)[len(basedir):]))
	} else {
		path = types.RelativePath(fulldir)
	}
	return path
}
