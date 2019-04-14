package apimodeler

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/labstack/echo"
)

// @TODO Move this to apimvc

type Context struct {
	Contexter
	RootDocumentor
	Controller ListController
}
type ContextArgs Context

func NewContext(args *ContextArgs) *Context {
	c := Context{}
	c = Context(*args)
	return &c
}

func (me *Context) GetRequestPath() (path types.UrlTemplate, sts status.Status) {
	for range only.Once {
		ec, sts := me.assertEchoContext()
		if is.Error(sts) {
			break
		}
		path = types.UrlTemplate(ec.Path())
	}
	return path, sts
}

func (me *Context) GetResponseBody() HttpResponseBody {
	return me.RootDocumentor.GetRootDocument()
}

func (me *Context) GetResponseStatus() (sc int, sts status.Status) {
	for range only.Once {
		ec, sts := me.assertEchoContext()
		if is.Error(sts) {
			break
		}
		sc = ec.Response().Status
	}
	return sc, sts
}

func (me *Context) SetResponseStatus(statuscode int) (sts status.Status) {
	for range only.Once {
		ec, sts := me.assertEchoContext()
		if is.Error(sts) {
			break
		}
		ec.Response().Status = statuscode
	}
	return sts
}

func (me *Context) SendResponse() (sts status.Status) {
	for range only.Once {
		ec, sts := me.assertEchoContext()
		if is.Error(sts) {
			break
		}
		err := ec.JSONPretty(
			ec.Response().Status,
			me.GetResponseBody(),
			"   ",
		)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("error while sending response for '%s'", ec.Path()),
			})
		}
	}
	return sts
}

func (me *Context) SetResponseHeader(name HttpHeaderName, value HttpHeaderValue) (sts status.Status) {
	for range only.Once {
		ec, sts := me.assertEchoContext()
		if is.Error(sts) {
			break
		}
		ec.Response().Header().Set(string(name), string(value))
	}
	return sts
}
func (me *Context) assertEchoContext() (ec echo.Context, sts status.Status) {
	for range only.Once {
		var ok bool
		ec, ok = me.Contexter.(echo.Context)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "context does not implement echo.Context",
			})
			break
		}
	}
	return ec, sts
}
