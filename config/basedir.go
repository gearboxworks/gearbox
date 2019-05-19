package config

import (
	"fmt"
	"gearbox/help"
	"gearbox/only"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/http"
	"path/filepath"
	"strings"
)

type BasedirMap map[types.Nickname]*Basedir

func (me BasedirMap) GetNickname(basedir types.AbsoluteDir, nickname ...types.Nickname) (nn types.Nickname) {
	var _nn types.Nickname
	if len(nickname) > 0 {
		_nn = nickname[0]
	}
	for range only.Once {
		var bd *Basedir
		for _nn, bd = range me {
			if bd.Basedir != basedir {
				continue
			}
			if _nn != "" && bd.Nickname == _nn {
				continue
			}
			nn = _nn
			break
		}
	}
	return nn
}

type Basedirs []*Basedir

type Basedir struct {
	Nickname types.Nickname    `json:"nickname"`
	Basedir  types.AbsoluteDir `json:"basedir"`
}

type BasedirArgs Basedir

func NewBasedir(nickname types.Nickname, basedir types.AbsoluteDir) *Basedir {
	return &Basedir{
		Nickname: nickname,
		Basedir:  basedir,
	}
}

func (me *Basedir) MaybeExpandBasedir() (sts Status) {
	for range only.Once {
		newdir, sts := util.MaybeExpandDir(me.Basedir)
		if is.Success(sts) {
			me.Basedir = newdir
		}
		sts = sts.SetMessage("for nickname '%s' %s",
			me.Basedir,
			sts.Message(),
		)
	}
	return sts
}

func (me *Basedir) Initialize() (sts Status) {
	for range only.Once {
		if me.Basedir == "" {
			sts = status.Fail(&status.Args{
				Message:    "Basedir.Basedir has no value",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("see %s", help.GetApiDocsUrl("basedirs")),
			})
			break
		}
		sts := me.MaybeExpandBasedir()
		if status.IsError(sts) {
			break
		}
		if me.Nickname == "" {
			me.Nickname = types.Nickname(filepath.Base(string(me.Basedir)))
		}
	}
	return sts
}

func (me BasedirMap) NicknameExists(nickname types.Nickname) (ok bool) {
	_, ok = me[nickname]
	return ok
}

func (me BasedirMap) BasedirExists(basedir types.AbsoluteDir) (ok bool) {
	for _, bd := range me {
		if bd.Basedir != basedir {
			continue
		}
		ok = true
		break
	}
	return ok
}

//func (me BasedirMap) BasedirDirExists(dir types.AbsoluteDir) (ok bool) {
//	for _, bd := range me {
//		if bd.Basedir != dir {
//			continue
//		}
//		ok = true
//		break
//	}
//	return ok
//}

func (me BasedirMap) GetBasedir(nickname types.Nickname) (bd *Basedir, sts Status) {
	bd, ok := me[nickname]
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("basedir '%s' not found", nickname),
		})
	}
	return bd, sts
}

func (me BasedirMap) DeleteBasedir(config Configer, nickname types.Nickname) (sts Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(nickname, &ValidateArgs{
			MustNotBeEmpty: true,
			MustNotEqual:   DefaultBasedirNickname,
			MustExist:      true,
			Config:         config,
		})
		if status.IsError(sts) {
			if strings.HasPrefix(sts.Message(), "nickname cannot equal") {
				sts = status.Wrap(sts, &status.Args{
					Message: fmt.Sprintf("cannot delete the base directory nicknamed '%s'",
						nickname,
					),
				})
			}
			break
		}
		var bd *Basedir
		bd, sts = me.GetBasedir(nickname)
		if is.Error(sts) {
			break
		}
		delete(me, nickname)
		sts = status.Success("basedir '%s' deleted",
			nickname,
		).SetDetail("'%s' was nickname for '%s'",
			nickname,
			bd.Basedir,
			/** Setting status code explicitly @see https://stackoverflow.com/a/2342589/102699 */
		).SetHttpStatus(http.StatusOK)
	}
	return sts
}

func (me BasedirMap) UpdateBasedir(config Configer, basedir *Basedir) (sts Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(basedir.Nickname, &ValidateArgs{
			Config:         config,
			MustNotBeEmpty: true,
			MustExist:      true,
		})
		if is.Error(sts) {
			break
		}
		ed, sts := util.MaybeExpandDir(basedir.Basedir)
		if is.Error(sts) {
			break
		}
		basedir.Basedir = ed
		sts = ValidateBasedir(basedir.Basedir, basedir.Nickname, &ValidateArgs{
			Config:         config,
			MustNotBeEmpty: true,
			MustBeOnDisk:   true,
			MustBeIn:       config.GetBasedirMap(),
			MustNotBeIn:    config.GetNicknameMap(),
			MustSucceed: func() (sts Status) {
				return me.ensureNonDuplicatedBasedir(basedir)
			},
		})
		if is.Error(sts) {
			break
		}
		sts = basedir.Initialize()
		if is.Error(sts) {
			break
		}
		me[basedir.Nickname] = basedir
		sts = status.Success("basedir '%s' updated", basedir.Nickname).
			SetDetail("'%s' is nickname for '%s'", basedir.Nickname, basedir.Basedir)
	}
	return sts
}

func (me BasedirMap) AddBasedir(config Configer, basedir *Basedir) (sts Status) {
	for range only.Once {
		sts = ValidateBasedirNickname(basedir.Nickname, &ValidateArgs{
			Config:         config,
			MustNotBeEmpty: true,
			MustNotExist:   true,
		})
		if is.Error(sts) {
			break
		}
		sts = ValidateBasedir(basedir.Basedir, basedir.Nickname, &ValidateArgs{
			Config:         config,
			MustNotBeEmpty: true,
			MustNotExist:   true,
			MustBeOnDisk:   true,
			MustNotBeIn:    config.GetBasedirMap(),
			MustSucceed: func() (sts Status) {
				return me.ensureNonDuplicatedBasedir(basedir)
			},
		})
		if is.Error(sts) {
			break
		}
		sts := basedir.Initialize()
		if is.Error(sts) {
			sts = sts.SetHttpStatus(http.StatusBadRequest)
			break
		}
		me[basedir.Nickname] = basedir
		sts = status.Success("base directory with nickname '%s' was added", basedir.Basedir).
			SetHttpStatus(http.StatusCreated).
			SetDetail("base directory with nickname '%s' is '%s'",
				basedir.Nickname,
				basedir.Basedir,
			)
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
		nn, ok := args.MustNotEqual.(string)
		if ok && nickname == types.Nickname(nn) {
			sts = status.YourBad("nickname cannot equal '%s'",
				nickname,
			)
			break
		}
		if args.MustNotBeEmpty && nickname == "" {
			sts = status.Fail(&status.Args{
				Message:    "basedir nickname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		nnExists := args.Config.NicknameExists(nickname)
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
				Message:    fmt.Sprintf("nickname '%s' already exists", nickname),
				HttpStatus: http.StatusConflict,
				ApiHelp:    apiHelp,
			})
			break
		}
		sts = status.Success("nickname '%s' validated", nickname)
	}
	return sts
}

func ValidateBasedir(basedir types.AbsoluteDir, nickname types.Nickname, args *ValidateArgs) (sts Status) {
	for range only.Once {
		sts = status.Success("base directory '%s' validated", basedir)
		if args.Config == nil {
			panic(fmt.Sprintf("Config property not passed in %T", args))
		}
		var apiHelp string
		if args.ApiHelpUrl != "" {
			apiHelp = fmt.Sprintf("see %s", args.ApiHelpUrl)
		}
		if args.MustNotBeEmpty && basedir == "" {
			sts = status.Fail(&status.Args{
				Message:    "base directory property 'basedir' is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		bdmap, ok := args.MustBeIn.(BasedirMap)
		if !ok {
			sts = status.Fail().
				SetMessage("unable to type assert `args.MustBeIn` to `BasedirMap` for basedir nickname '%s'.",
					nickname,
				)
			break
		}
		_, ok = bdmap[nickname]
		if !ok {
			sts = status.Fail(&status.Args{
				ApiHelp:    apiHelp,
				HttpStatus: http.StatusBadRequest,
				Message:    fmt.Sprintf("nickname for base directory '%s' not found", basedir),
			})
			break
		}
		var nnm NicknameMap
		nnm, ok = args.MustNotBeIn.(NicknameMap)
		if !ok {
			sts = status.Fail().
				SetMessage("unable to type assert `args.MustNotBeIn` to `NicknameMap` for basedir nickname '%s'.",
					nickname,
				)
			break
		}
		var nn types.Nickname
		nn, ok = nnm[basedir]
		if ok && args.IgnoreCurrent && nn != nickname {
			sts = status.Fail().
				SetMessage("base directory '%s' already exists as nickname '%s'", basedir, nn)
		}
		if ok && !args.IgnoreCurrent {
			sts = status.Fail().
				SetMessage("base directory '%s' already exists", basedir, nn)
		}
		if is.Error(sts) {
			sts = sts.
				SetHelp(status.ApiHelp, apiHelp).
				SetHttpStatus(http.StatusBadRequest)
			break
		}
		if !ok && args.MustExist {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("base directory '%s' not found", basedir),
				HttpStatus: http.StatusNotFound,
				ApiHelp:    apiHelp,
			})
			break
		}
		sts = args.MustSucceed()
		if is.Error(sts) {
			break
		}
		if !args.MustBeOnDisk && !args.MustNotBeOnDisk {
			break
		}
		bdOnDisk := util.DirExists(basedir)
		if args.MustBeOnDisk && !bdOnDisk {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("base directory '%s' does not exist", basedir),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		if args.MustNotBeOnDisk && bdOnDisk {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("base directory '%s' already exists", basedir),
				HttpStatus: http.StatusConflict,
				ApiHelp:    apiHelp,
			})
			break
		}
	}
	return sts
}

func (me BasedirMap) ensureNonDuplicatedBasedir(bd *Basedir) (sts Status) {
	nn := me.GetNickname(bd.Basedir)
	if nn != "" && nn != bd.Nickname {
		sts = status.Fail().SetMessage("base directory '%s' already has nickname '%s'",
			bd.Basedir,
			me.GetNickname(bd.Basedir),
		)
	}
	return sts
}
