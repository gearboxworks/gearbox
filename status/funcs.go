package status

import (
	"errors"
	"fmt"
	"gearbox/only"
	"net/http"
)

func HttpStatus(err error) int {
	sts, ok := err.(Status)
	if !ok {
		return 0
	}
	return sts.HttpStatus()
}

func Message(err error) string {
	sts, ok := err.(Status)
	if !ok {
		return ""
	}
	return sts.Message()
}

func Cause(err error) error {
	sts, ok := err.(Status)
	if !ok {
		return err
	}
	return sts.Cause()
}

func Help(err error, helptype ...HelpType) string {
	sts, ok := err.(Status)
	if !ok {
		return ""
	}
	if len(helptype) == 0 {
		helptype = []HelpType{AllHelp}
	}
	return sts.GetHelp(helptype[0])
}

func IsError(err error) bool {
	sts, ok := err.(Status)
	if !ok {
		return err != nil
	}
	return sts.IsError()
}

func IsSuccess(err error) bool {
	sts, ok := err.(Status)
	if !ok {
		return err == nil
	}
	return !sts.IsError() || sts.HttpStatus() == 0
}

func Wrap(err error, args ...*Args) Status {
	var _args *Args
	if len(args) == 0 || args[0] == nil {
		_args = &Args{}
	} else {
		_args = args[0]
	}
	_args.Cause = err
	if _args.Message == "" {
		_args.Message = err.Error()
	}
	return NewStatus(_args)
}

func SimpleSuccess() Status {
	return Success("everything a-ok")
}

func Success(msg string, args ...interface{}) Status {
	return NewStatus(&Args{
		Success:    true,
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusOK,
	})
}

func Fail(args *Args) Status {
	args.Success = false
	return NewStatus(args)
}

func YourBad(msg string, args ...interface{}) Status {
	return Fail(&Args{
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusBadRequest,
	})
}

func OurBad(msg string, args ...interface{}) Status {
	return Fail(&Args{
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusInternalServerError,
	})
}

func NewStatus(args *Args) (sts *S) {
	for range only.Once {
		sts = &S{
			success:    args.Success,
			message:    args.Message,
			httpstatus: args.HttpStatus,
			cause:      args.Cause,
			data:       args.Data,
			help: HelpTypeMap{
				AllHelp: &args.Help,
				ApiHelp: &args.ApiHelp,
				CliHelp: &args.CliHelp,
			},
		}

		if sts.message == "" && sts.cause != nil {
			sts.message = sts.cause.Error()
		}

		if !sts.success && sts.cause == nil {
			sts.cause = errors.New(sts.message)
		}

		if sts.httpstatus == 0 {
			sts.httpstatus = http.StatusInternalServerError
		}

		if *sts.help[AllHelp] == "" {
			help := ContactSupportHelp()
			sts.help[AllHelp] = &help
		}

		if *sts.help[ApiHelp] == "" {
			sts.help[ApiHelp] = sts.help[AllHelp]
		}

		if *sts.help[CliHelp] == "" {
			sts.help[CliHelp] = sts.help[AllHelp]
		}

	}
	return sts
}

func ContactSupportHelp() string {
	return "contact support"
}
