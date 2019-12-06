// A simple wrapper around osbridge.OsBridger.
// This makes it much easier to separate the EventBroker code into it's own package later on.
package osdirs

import (
	"fmt"
	"os"
)

type Paths []Path
type Path struct {
	Dir  Dir
	File File
}

func NewPaths() *Paths {
	return &Paths{}
}

func (me *Paths) AppendFile(file ...string) *Paths {
	var ret Paths
	if me != nil {
		ret = *me
	}
	for _, s := range file {
		ret = append(ret, *Split(s))
	}
	return &ret
}

func (me *Paths) AppendDir(dir ...string) *Paths {
	var ret Paths
	if me != nil {
		ret = *me
	}
	for _, s := range dir {
		if s == "" {
			continue
		}
		ret = append(ret, Path{Dir: s})
	}
	return &ret
}

func (me *Paths) CreateIfNotExists() (err error) {
	for _, p := range *me {
		if p.Dir == "" {
			continue
		}
		_, err = CreateIfNotExists(p.Dir)
		if err != nil {
			break
		}
	}
	return err
}

func (me *Path) CreateIfNotExists() (created bool, err error) {
	created, err = CreateIfNotExists(me.Dir)
	if err != nil {
		fmt.Printf("CreateFileIfNotExists PATH: '%s'\n", me.String())
		err = os.MkdirAll(me.Dir, os.ModePerm)
		created = true
	}
	return created, err
}

func (me *Path) String() string {
	return fmt.Sprintf("%s%c%s",
		me.Dir,
		os.PathSeparator,
		me.File,
	)
}
