package gearbox

import (
	"errors"
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"path/filepath"
	"strings"
)

const PrimaryBasedirNickname = "primary"

type BasedirMap map[string]*Basedir

type Basedirs []*Basedir

type Basedir struct {
	Nickname string `json:"nickname"`
	HostDir  string `json:"host_dir"`
	BoxDir   string `json:"box_dir"`
	Error    error  `json:"-"`
}
type BasedirArgs Basedir

func NewBasedir(hostDir string, args ...*BasedirArgs) *Basedir {
	var _args *BasedirArgs
	if len(args) == 0 {
		_args = &BasedirArgs{HostDir: hostDir}
	} else {
		_args = args[0]
	}
	bd := Basedir(*_args)
	bd.HostDir = hostDir
	bd.Initialize()
	return &bd
}

func ExpandHostBasedirPath(gb *Gearbox, nickname string, path string) (fp string, status stat.Status) {
	status = gb.ValidateBasedirNickname(nickname, &validateArgs{
		MustNotBeEmpty: true,
		MustExist:      true,
		ApiHelpUrl:     api.GetApiDocsUrl(gb.RequestType),
	})
	if !status.IsError() {
		bd, _ := gb.Config.GetHostBasedir(nickname)
		fp = filepath.FromSlash(fmt.Sprintf("%s/%s", bd, path))
	}
	return fp, status
}

func (me *Basedir) MaybeExpandDir() (status stat.Status) {
	for range only.Once {
		origDir := me.HostDir
		if strings.HasPrefix(me.HostDir, "~") {
			dir, err := homedir.Expand(me.HostDir)
			if err != nil {
				status = stat.NewStatus(&stat.Args{
					Error:      err,
					HttpStatus: http.StatusInternalServerError,
					Message: fmt.Sprintf("could not expand dir '%s' for '%s'",
						me.HostDir,
						me.Nickname,
					),
				})
				break
			}
			me.HostDir = dir
		}
		status = stat.NewOkStatus("directory expanded from '%s' to '%s'",
			origDir,
			me.HostDir,
		)
	}
	return status
}

func (me *Basedir) Initialize() (status stat.Status) {
	for range only.Once {
		if me.HostDir == "" {
			me.Error = errors.New("Basedir.HostDir has no value")
			status = stat.NewStatus(&stat.Args{
				Error:      me.Error,
				Message:    me.Error.Error(),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("see %s", api.GetApiDocsUrl("basedirs")),
			})
			break
		}
		status := me.MaybeExpandDir()
		if status.IsError() {
			me.Error = status.Error
			break
		}
		if me.Nickname == "" {
			me.Nickname = filepath.Base(me.HostDir)
		}
		if me.BoxDir == PrimaryBasedirNickname {
			me.BoxDir = boxBasedir
			break
		}
		if me.BoxDir == "" || me.BoxDir == boxBasedir {
			me.BoxDir = filepath.FromSlash(fmt.Sprintf("%s/%s", boxBasedir, me.Nickname))
		}
	}
	return status
}

func (me BasedirMap) NamedBasedirExists(nickname string) (ok bool) {
	_, ok = me[nickname]
	return ok
}

func (me BasedirMap) BasedirExists(dir string) (ok bool) {
	for _, bd := range me {
		if bd.HostDir != dir {
			continue
		}
		ok = true
		break
	}
	return ok
}

func (me BasedirMap) GetNamedBasedir(nickname string) *Basedir {
	bd, _ := me[nickname]
	return bd
}

func (me BasedirMap) DeleteNamedBasedir(gb *Gearbox, nickname string) (status stat.Status) {
	for range only.Once {
		status = gb.ValidateBasedirNickname(nickname, &validateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
			ApiHelpUrl:     api.GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := me.GetNamedBasedir(nickname)
		delete(me, nickname)
		status = stat.NewSuccessStatus(
			http.StatusOK,
			"named base dir '%s' ('%s') deleted",
			nickname,
			bd.HostDir,
		)
	}
	return status
}

func (me BasedirMap) UpdateBasedir(gb *Gearbox, nickname string, dir string) (status stat.Status) {
	for range only.Once {
		status = gb.ValidateBasedirNickname(nickname, &validateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
			ApiHelpUrl:     api.GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := me.GetNamedBasedir(nickname)
		bd.HostDir = dir
		status = bd.MaybeExpandDir()
		if status.IsError() {
			break
		}
		bd.Initialize()
		if status.IsError() {
			break
		}
		status = stat.NewSuccessStatus(
			http.StatusOK,
			"named base dir '%s' updated to: '%s'",
			nickname,
			bd.HostDir,
		)
	}
	return status
}

func (me BasedirMap) AddBasedir(gb *Gearbox, dir string, nickname ...string) (status stat.Status) {
	for range only.Once {
		var nn string
		if len(nickname) > 0 {
			nn = nickname[0]
		}
		status = gb.ValidateBasedirNickname(nn, &validateArgs{
			MustNotBeEmpty: true,
			MustNotExist:   true,
			ApiHelpUrl:     api.GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := NewBasedir(dir, &BasedirArgs{
			BoxDir:   gb.Config.BoxBasedir,
			Nickname: nn,
		})
		if bd.Error != nil {
			status = stat.NewStatus(&stat.Args{
				HttpStatus: http.StatusBadRequest,
				Error:      bd.Error,
			})
			if dir == "" {
				status.Message = fmt.Sprint("value provide for base dir in 'host_dir' property was empty")
			} else {
				status.Message = fmt.Sprintf("could add add base dir '%s'; the ~ could not be expanded", dir)
			}
			break
		}
		me[bd.Nickname] = bd
		status = stat.NewSuccessStatus(
			http.StatusCreated,
			"base dir '%s' added",
			bd.HostDir,
		)
	}
	return status
}

func ValidateBasedirNickname(nickname string, args *validateArgs) (status stat.Status) {
	for range only.Once {
		var apiHelp string
		if args.ApiHelpUrl != "" {
			apiHelp = fmt.Sprintf("see %s", args.ApiHelpUrl)
		}
		if args.MustNotBeEmpty && nickname == "" {
			status = stat.NewStatus(&stat.Args{
				Failed:     true,
				Message:    "basedir nickname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		nnExists := args.Gearbox.NamedBasedirExists(nickname)
		if args.MustExist && !nnExists {
			status = stat.NewStatus(&stat.Args{
				Failed:     true,
				Message:    fmt.Sprintf("nickname '%s' does not exist", nickname),
				HttpStatus: http.StatusNotFound,
				ApiHelp:    apiHelp,
			})
			break
		}
		if args.MustNotExist && nnExists {
			status = stat.NewStatus(&stat.Args{
				Failed:     true,
				Message:    fmt.Sprintf("nickname '%s' already exists", nickname),
				HttpStatus: http.StatusInternalServerError,
				ApiHelp:    apiHelp,
			})
			break
		}
		status = stat.NewOkStatus("nickname '%s' validated", nickname)
	}
	if status.IsError() {
		status.Error = errors.New(status.Message)
	}
	return status
}
