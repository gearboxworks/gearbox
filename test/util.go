package test

import (
	"strings"
)

func ParseRelativePath(basedir, fulldir string) (path string) {
	if strings.HasPrefix(fulldir, basedir) {
		path = string([]byte(fulldir)[len(basedir):])
	} else {
		path = fulldir
	}
	return path
}
