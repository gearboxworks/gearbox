package config

import (
	"errors"
	"fmt"
	"gearbox/help"
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
	Basedir  types.AbsoluteDir `json:"basedir"`
	Err      error             `json:"-"`
}
type BasedirArgs Basedir

func NewBasedir(nickname types.Nickname, basedir types.AbsoluteDir) *Basedir {
	return &Basedir{
		Nickname: nickname,
		Basedir:  basedir,
	}
}
func (me *Basedir) MaybeExpandHostDir() (sts Status) {
	for range only.Once {
		origDir := me.Basedir
		if strings.HasPrefix(string(me.Basedir), "~") {
			dir, err := homedir.Expand(string(me.Basedir))
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("could not expand '%s' dir '%s'",
						me.Nickname,
						me.Basedir,
					),
				})
				break
			}
			me.Basedir = types.AbsoluteDir(dir)
		}
		sts = status.Success("directory '%s' expanded from '%s' to '%s'",
			me.Nickname,
			origDir,
			me.Basedir,
		)
	}
	return sts
}

func (me *Basedir) Initialize() (sts Status) {
	for range only.Once {
		if me.Basedir == "" {
			me.Err = errors.New("Basedir.Basedir has no value")
			sts = status.Wrap(me.Err, &status.Args{
				Message:    me.Err.Error(),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("see %s", help.GetApiDocsUrl("basedirs")),
			})
			break
		}
		sts := me.MaybeExpandHostDir()
		if status.IsError(sts) {
			me.Err = sts.Cause()
			break
		}
		if me.Nickname == "" {
			me.Nickname = types.Nickname(filepath.Base(string(me.Basedir)))
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
		if bd.Basedir != dir {
			continue
		}
		ok = true
		break
	}
	return ok
}

func (me BasedirMap) GetNamedBasedir(nickname types.Nickname) (bd *Basedir, sts Status) {
	bd, ok := me[nickname]
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("basedir '%s' not found", nickname),
		})
	}
	return bd, sts
}

func (me BasedirMap) DeleteNamedBasedir(nickname types.Nickname) (sts Status) {
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
			bd.Basedir,
		)
	}
	return sts
}

func (me BasedirMap) UpdateBasedir(nickname types.Nickname, dir types.AbsoluteDir) (sts Status) {
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
		bd.Basedir = dir
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
			bd.Basedir,
		)
	}
	return sts
}

func (me BasedirMap) AddBasedir(config Configer, basedir *Basedir) (sts Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(basedir.Nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustNotExist:   true,
			Config:         config,
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
		sts = status.Success("base dir '%s' added", basedir.Basedir)
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func ValidateBasedirNickname(nickname types.Nickname, args *ValidateArgs) (sts Status) {
	for range only.Once {
		if args.Config == nil {
			panic(fmt.Sprintf("Config property not passed in %T", args))
		}
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
