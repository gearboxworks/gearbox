package config

import (
	"errors"
	"fmt"
	"gearbox/api"
	"gearbox/box"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"path/filepath"
	"strings"
)

type BasedirMap map[types.Nickname]*Basedir

type Basedirs []*Basedir

type Basedir struct {
	Nickname types.Nickname    `json:"nickname"`
	HostDir  types.AbsoluteDir `json:"host_dir"`
	BoxDir   types.AbsoluteDir `json:"box_dir"`
	Err      error             `json:"-"`
}
type BasedirArgs Basedir

func NewBasedir(hostdir types.AbsoluteDir, args ...*BasedirArgs) *Basedir {
	var _args *BasedirArgs
	if len(args) == 0 {
		_args = &BasedirArgs{HostDir: hostdir}
	} else {
		_args = args[0]
	}
	bd := Basedir(*_args)
	bd.HostDir = hostdir
	return &bd
}
func (me *Basedir) MaybeExpandHostDir() (sts status.Status) {
	for range only.Once {
		origDir := me.HostDir
		if strings.HasPrefix(string(me.HostDir), "~") {
			dir, err := homedir.Expand(string(me.HostDir))
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("could not expand '%s' dir '%s'",
						me.Nickname,
						me.HostDir,
					),
				})
				break
			}
			me.HostDir = types.AbsoluteDir(dir)
		}
		sts = status.Success("directory '%s' expanded from '%s' to '%s'",
			me.Nickname,
			origDir,
			me.HostDir,
		)
	}
	return sts
}

func (me *Basedir) Initialize() (sts status.Status) {
	for range only.Once {
		if me.HostDir == "" {
			me.Err = errors.New("Basedir.HostDir has no value")
			sts = status.Wrap(me.Err, &status.Args{
				Message:    me.Err.Error(),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("see %s", api.GetApiDocsUrl("basedirs")),
			})
			break
		}
		sts := me.MaybeExpandHostDir()
		if status.IsError(sts) {
			me.Err = sts.Cause()
			break
		}
		if me.Nickname == "" {
			me.Nickname = types.Nickname(filepath.Base(string(me.HostDir)))
		}
		if me.Nickname == PrimaryBasedirNickname {
			me.BoxDir = box.Basedir
			break
		}
		if me.BoxDir == "" || me.BoxDir == box.Basedir {
			me.BoxDir = types.AbsoluteDir(filepath.FromSlash(fmt.Sprintf("%s/%s",
				box.Basedir,
				me.Nickname,
			)))
		}
	}
	return sts
}

func (me BasedirMap) NamedBasedirExists(nickname types.Nickname) (ok bool) {
	_, ok = me[nickname]
	return ok
}

func (me BasedirMap) BasedirExists(dir types.AbsoluteDir) (ok bool) {
	for _, bd := range me {
		if bd.HostDir != dir {
			continue
		}
		ok = true
		break
	}
	return ok
}

func (me BasedirMap) GetNamedBasedir(nickname types.Nickname) (bd *Basedir, sts status.Status) {
	bd, ok := me[nickname]
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("basedir '%s' not found", nickname),
		})
	}
	return bd, sts
}

func (me BasedirMap) DeleteNamedBasedir(nickname types.Nickname) (sts status.Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
		})
		if status.IsError(sts) {
			break
		}
		bd, sts := me.GetNamedBasedir(nickname)
		if is.Error(sts) {
			break
		}
		delete(me, nickname)
		sts = status.Success("named base dir '%s' ('%s') deleted",
			nickname,
			bd.HostDir,
		)
	}
	return sts
}

func (me BasedirMap) UpdateBasedir(nickname types.Nickname, dir types.AbsoluteDir) (sts status.Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
		})
		if status.IsError(sts) {
			break
		}
		bd, sts := me.GetNamedBasedir(nickname)
		if is.Error(sts) {
			break
		}
		bd.HostDir = dir
		sts = bd.MaybeExpandHostDir()
		if status.IsError(sts) {
			break
		}
		sts = bd.Initialize()
		if status.IsError(sts) {
			break
		}
		sts = status.Success(
			"named base dir '%s' updated to: '%s'",
			nickname,
			bd.HostDir,
		)
	}
	return sts
}

func (me BasedirMap) AddBasedir(basedir *Basedir) (sts status.Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(basedir.Nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustNotExist:   true,
		})
		if is.Error(sts) {
			break
		}
		sts := basedir.Initialize()
		if is.Error(sts) {
			sts.SetHttpStatus(http.StatusBadRequest)
			break
		}
		me[basedir.Nickname] = basedir
		sts = status.Success("base dir '%s' added", basedir.HostDir)
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func ValidateBasedirNickname(nickname types.Nickname, args *ValidateArgs) (sts status.Status) {
	for range only.Once {
		var apiHelp string
		if args.ApiHelpUrl != "" {
			apiHelp = fmt.Sprintf("see %s", args.ApiHelpUrl)
		}
		if args.MustNotBeEmpty && nickname == "" {
			sts = status.Fail(&status.Args{
				Message:    "basedir nickname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		nnExists := args.Config.NamedBasedirExists(nickname)
		if args.MustExist && !nnExists {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("nickname '%s' does not exist", nickname),
				HttpStatus: http.StatusNotFound,
				ApiHelp:    apiHelp,
			})
			break
		}
		if args.MustNotExist && nnExists {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("nickname '%s' already exists", nickname),
				ApiHelp: apiHelp,
			})
			break
		}
		sts = status.Success("nickname '%s' validated", nickname)
	}
	return sts
}
