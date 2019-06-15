// A simple wrapper around oss.OsSupporter.
// This makes it much easier to separate the EventBroker code into it's own package later on.
package ospaths

import (
	"fmt"
	"gearbox/eventbroker/only"
	oss "gearbox/os_support"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	DefaultBaseDir = "dist/eventbroker"

	defaultLogBaseDir = "logs"
	defaultEtcBaseDir = "etc"
)

type Path struct {
	Dir
	File
}
type Paths []Path

type Dir  string
type Dirs []Dir

type File string
type Files []File

type BasePaths struct {
	UserHomeDir           Dir
	SuggestedBasedir      Dir
	UserConfigDir         Dir
	AdminRootDir          Dir
	CacheDir              Dir
	EventBrokerDir        Dir
	EventBrokerWorkingDir Dir
	EventBrokerLogDir     Dir
	EventBrokerEtcDir     Dir
	LocalDir              Dir
	mutex                 sync.RWMutex
}
//type OsSupport     oss.OsSupporter


func New(subdir string) *BasePaths {

	var ret BasePaths

	ret.UserHomeDir = Dir(oss.Get().GetUserHomeDir())
	ret.SuggestedBasedir = Dir(oss.Get().GetSuggestedBasedir())
	ret.UserConfigDir = Dir(oss.Get().GetUserConfigDir())
	ret.AdminRootDir = Dir(oss.Get().GetAdminRootDir())
	ret.CacheDir = Dir(oss.Get().GetCacheDir())

	ret.LocalDir = Dir(filepath.FromSlash("/usr/local"))
	ret.EventBrokerDir = *ret.UserConfigDir.AddToPath(DefaultBaseDir)
	ret.EventBrokerLogDir = *ret.UserConfigDir.AddToPath(DefaultBaseDir, defaultLogBaseDir)
	ret.EventBrokerEtcDir = *ret.UserConfigDir.AddToPath(DefaultBaseDir, defaultEtcBaseDir)
	ret.EventBrokerWorkingDir = ret.EventBrokerDir
	//ret.EventBrokerDir = Dir(filepath.FromSlash(fmt.Sprintf("%s/dist/eventbroker", ret.UserConfigDir)))

	//ret.ChannelsDir = Dir(filepath.FromSlash(fmt.Sprintf("%s", ret.EventBrokerDir)))
	//ret.MqttClientDir = Dir(filepath.FromSlash(fmt.Sprintf("%s", ret.EventBrokerDir)))

	return &ret
}


func (me *Dir) AddToPath(dir ...string) *Dir {

	var ret Dir
	var d []string

	d = append(d, string(*me))
	d = append(d, dir...)

	ret = Dir(filepath.FromSlash(strings.Join(d, "/")))

	return &ret
}


func (me *Dir) AddFileToPath(format string, fn ...interface{}) *Dir {

	var ret Dir
	var d []string

	d = append(d, string(*me))
	d = append(d, fmt.Sprintf(format, fn...))

	ret = Dir(filepath.FromSlash(strings.Join(d, "/")))

	return &ret
}


func (me *Dir) DirExists() error {

	var err error

	if _, err = os.Stat(me.String()); os.IsNotExist(err) {
		//fmt.Printf("Not exists PATH: '%s'\n", me.String())
	}

	return err
}


func (me *Dir) CreateIfNotExists() (created bool, err error) {

	if me.DirExists() != nil {
		//fmt.Printf("CreateDirIfNotExists PATH: '%s'\n", me.String())
		err = os.MkdirAll(me.String(), os.ModePerm)
		if err == nil {
			created = true
		}
	}

	return created, err
}


func (me *Dirs) Append(dir ...string) *Dirs {

	var ret Dirs
	if me != nil {
		ret = *me
	}

	for _, s := range dir {
		ret = append(ret, Dir(s))
	}

	return &ret
}


func NewPath() *Paths {

	var ret Paths

	return &ret
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

		ret = append(ret, Path{Dir: Dir(s), File: ""})
	}

	return &ret
}


func (me *BasePaths) CreateIfNotExists() (err error) {

	for range only.Once {
		_, err = me.EventBrokerDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = me.EventBrokerEtcDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = me.EventBrokerLogDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = me.EventBrokerWorkingDir.CreateIfNotExists()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Paths) CreateIfNotExists() (err error) {

	for _, p := range *me {
		if p.String() == "" {
			continue
		}

		_, err = p.Dir.CreateIfNotExists()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Path) CreateIfNotExists() (created bool, err error) {

	created, err = me.Dir.CreateIfNotExists()
	if err != nil {
		fmt.Printf("CreateFileIfNotExists PATH: '%s'\n", me.String())
		err = os.MkdirAll(me.String(), os.ModePerm)
		created = true
	}

	return created, err
}


//func (me *Path) DirName() (created bool, err error) {
//
//	created, err = me.CreateIfNotExists()
//	if err != nil {
//		fmt.Printf("CreateFileIfNotExists PATH: '%s'\n", me.String())
//		err = os.MkdirAll(me.String(), os.ModePerm)
//		created = true
//	}
//
//	return created, err
//}


func (me *Dir) String() string {

	return string(*me)
}


func Split(fn string) *Path {

	var pn Path

	d, f := filepath.Split(fn)
	pn.Dir = Dir(d)
	pn.File = File(f)

	return &pn
}

