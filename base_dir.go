package gearbox

import (
	"errors"
	"fmt"
	"gearbox/only"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"path/filepath"
	"strings"
)

const PrimaryBaseDirNickname = "primary"

type BaseDirMap map[string]*BaseDir

type BaseDirs []*BaseDir

type BaseDir struct {
	Nickname string `json:"nickname"`
	HostDir  string `json:"host_dir"`
	VmDir    string `json:"vm_dir"`
	Error    error  `json:"-"`
}
type BaseDirArgs BaseDir

func NewBaseDir(hostDir string, args ...*BaseDirArgs) *BaseDir {
	var _args *BaseDirArgs
	if len(args) == 0 {
		_args = &BaseDirArgs{HostDir: hostDir}
	} else {
		_args = args[0]
	}
	bd := BaseDir(*_args)
	bd.HostDir = hostDir
	bd.Initialize()
	return &bd
}

func ExpandBaseDirPath(gb *Gearbox, nickname string, path string) (fp string, status Status) {
	status = gb.ValidateBaseDirNickname(nickname, &validateArgs{
		MustNotBeEmpty: true,
		MustExist:      true,
		ApiHelpUrl:     GetApiDocsUrl(gb.RequestType),
	})
	if !status.IsError() {
		fp = fmt.Sprintf("%s/%s",
			gb.Config.GetHostBaseDir(nickname),
			path,
		)
	}
	return fp, status
}

func (me *BaseDir) MaybeExpandDir() (status Status) {
	for range only.Once {
		origDir := me.HostDir
		if strings.HasPrefix(me.HostDir, "~") {
			dir, err := homedir.Expand(me.HostDir)
			if err != nil {
				status = NewStatus(&StatusArgs{
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
		status = NewOkStatus("directory expanded from '%s' to '%s'",
			origDir,
			me.HostDir,
		)
	}
	return status
}

func (me *BaseDir) Initialize() (status Status) {
	for range only.Once {
		if me.HostDir == "" {
			me.Error = errors.New("BaseDir.HostDir has no value")
			status = NewStatus(&StatusArgs{
				Error:      me.Error,
				Message:    me.Error.Error(),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("see %s", GetApiDocsUrl("basedirs")),
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
		if me.VmDir == PrimaryBaseDirNickname {
			me.VmDir = vmBaseDir
			break
		}
		if me.VmDir == "" || me.VmDir == vmBaseDir {
			me.VmDir = fmt.Sprintf("%s/%s", vmBaseDir, me.Nickname)
		}
	}
	return status
}

func (me BaseDirMap) NamedBaseDirExists(nickname string) (ok bool) {
	_, ok = me[nickname]
	return ok
}

func (me BaseDirMap) BaseDirExists(dir string) (ok bool) {
	for _, bd := range me {
		if bd.HostDir != dir {
			continue
		}
		ok = true
		break
	}
	return ok
}

func (me BaseDirMap) GetNamedBaseDir(nickname string) *BaseDir {
	bd, _ := me[nickname]
	return bd
}

func (me BaseDirMap) DeleteNamedBaseDir(gb *Gearbox, nickname string) (status Status) {
	for range only.Once {
		status = gb.ValidateBaseDirNickname(nickname, &validateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
			ApiHelpUrl:     GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := me.GetNamedBaseDir(nickname)
		delete(me, nickname)
		status = NewSuccessStatus(
			http.StatusOK,
			"named base dir '%s' ('%s') deleted",
			nickname,
			bd.HostDir,
		)
	}
	return status
}

func (me BaseDirMap) UpdateBaseDir(gb *Gearbox, nickname string, dir string) (status Status) {
	for range only.Once {
		status = gb.ValidateBaseDirNickname(nickname, &validateArgs{
			MustNotBeEmpty: true,
			MustExist:      true,
			ApiHelpUrl:     GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := me.GetNamedBaseDir(nickname)
		bd.HostDir = dir
		status = bd.MaybeExpandDir()
		if status.IsError() {
			break
		}
		bd.Initialize()
		if status.IsError() {
			break
		}
		status = NewSuccessStatus(
			http.StatusOK,
			"named base dir '%s' updated to: '%s'",
			nickname,
			bd.HostDir,
		)
	}
	return status
}

func (me BaseDirMap) AddBaseDir(gb *Gearbox, dir string, nickname ...string) (status Status) {
	for range only.Once {
		var nn string
		if len(nickname) > 0 {
			nn = nickname[0]
		}
		status = gb.ValidateBaseDirNickname(nn, &validateArgs{
			MustNotBeEmpty: true,
			MustNotExist:   true,
			ApiHelpUrl:     GetApiDocsUrl(gb.RequestType),
		})
		if status.IsError() {
			break
		}
		bd := NewBaseDir(dir, &BaseDirArgs{
			VmDir:    gb.Config.VmBaseDir,
			Nickname: nn,
		})
		if bd.Error != nil {
			status = NewStatus(&StatusArgs{
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
		status = NewSuccessStatus(
			http.StatusCreated,
			"base dir '%s' added",
			bd.HostDir,
		)
	}
	return status
}

func ValidateBaseDirNickname(nickname string, args *validateArgs) (status Status) {
	for range only.Once {
		var apiHelp string
		if args.ApiHelpUrl != "" {
			apiHelp = fmt.Sprintf("see %s", args.ApiHelpUrl)
		}
		if args.MustNotBeEmpty && nickname == "" {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    "basedir nickname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		nnExists := args.Gearbox.NamedBaseDirExists(nickname)
		if args.MustExist && !nnExists {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    fmt.Sprintf("nickname '%s' does not exist", nickname),
				HttpStatus: http.StatusNotFound,
				ApiHelp:    apiHelp,
			})
			break
		}
		if args.MustNotExist && nnExists {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    fmt.Sprintf("nickname '%s' already exists", nickname),
				HttpStatus: http.StatusInternalServerError,
				ApiHelp:    apiHelp,
			})
			break
		}
		status = NewOkStatus("nickname '%s' validated", nickname)
	}
	return status
}
