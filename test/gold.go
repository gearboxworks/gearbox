package test

import (
	"flag"
	"fmt"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/only"
	"io/ioutil"
	"os"
)

const GoldFilepathTemplate = "fixtures/gold/%s.gold.json"

var updateGold = flag.Bool("update-gold", false, "update gold files")

type GoldFile struct {
	filepath string
	Status   status.Status
}

func NewGoldFile(fileid string) *GoldFile {
	gf := GoldFile{}
	for range only.Once {
		dir, err := os.Getwd()
		if err != nil {
			gf.Status = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("cannot get working directory for '%s'", fileid),
			})
			break
		}
		absdir := util.ParentDir(types.Dir(dir))
		gf.filepath = fmt.Sprintf("%s/%s", absdir, fmt.Sprintf(GoldFilepathTemplate, fileid))
	}
	return &gf
}

func (me *GoldFile) DoUpdate() bool {
	return *updateGold
}

func (me *GoldFile) DoRunTest() bool {
	return !*updateGold
}

func (me *GoldFile) Read() (b []byte, sts status.Status) {
	for range only.Once {
		if status.IsError(me.Status) {
			break
		}
		var err error
		b, err = ioutil.ReadFile(me.filepath)
		if err != err {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("cannot open/read '%s'", me.filepath),
			})
			me.Status = sts
			break
		}
	}
	return b, sts
}

func (me *GoldFile) Write(gold []byte) (sts status.Status) {
	for range only.Once {
		if status.IsError(me.Status) {
			break
		}
		if !*updateGold {
			break
		}
		dir := util.FileDir(types.Filepath(me.filepath))
		_ = os.MkdirAll(string(dir), os.ModePerm)
		err := ioutil.WriteFile(me.filepath, gold, os.ModePerm)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("cannot write to '%s'", me.filepath),
			})
			me.Status = sts
			break
		}
	}
	return sts
}
