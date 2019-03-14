package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

type RequestContext struct {
	*Api
	ResourceName
	echo.Context
}

func NewRequestContext(_api *Api, resourceName ResourceName) *RequestContext {
	return &RequestContext{
		Api:          _api,
		ResourceName: resourceName,
	}
}
func (me *RequestContext) Param(name string) (value string) {
	if me.Context != nil {
		value = me.Context.Param(name)
	}
	return value
}

func (me *RequestContext) WrapHandler(next HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		me.Context = ctx
		return next(me)
	}
}

func (me *RequestContext) ReadRequestBody() ([]byte, error) {
	return ioutil.ReadAll(me.Context.Request().Body)
}

func (me *RequestContext) CloseRequestBody() {
	_ = me.Context.Request().Body.Close()
}

func (me *RequestContext) UnmarshalFromRequest(obj interface{}) (err error) {
	for range only.Once {
		apiHelp := GetApiHelp(me.ResourceName)
		defer me.CloseRequestBody()
		b, err := me.ReadRequestBody()
		if err != nil {
			err = util.AddHelpToError(
				errors.New("could not read request body"),
				apiHelp,
			)
			break
		}
		err = json.Unmarshal(b, &obj)
		if err != nil {
			err = util.AddHelpToError(
				fmt.Errorf("unexpected format for request body: '%s'", string(b)),
				apiHelp,
			)
			break
		}
	}
	return err
}

// @TODO Add ?format=yes to pretty print JSON
func (rc *RequestContext) JsonMarshalHandler(js interface{}) (status *Status) {
	var err error
	ctx := rc.Context
	for range only.Once {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		var ok bool
		success := true
		status, ok = js.(*Status)
		if ok {
			success = false
			js = status.ToResponse()
			ctx.Response().Status = status.StatusCode
		}
		httpStatus := ctx.Response().Status
		r := *rc.Api.Defaults.Clone()
		if rdg, ok := js.(ResponseDataGetter); !ok {
			r.Data = js
		} else {
			r.Data = rdg.(ResponseDataGetter).GetResponseData()
		}
		path := ctx.Path()
		if path != "/" {
			r.Links = make(Links, 0)
		}
		if slg, ok := js.(SelfLinkGetter); ok {
			r.Links[SelfResource] = slg.GetApiSelfLink()
		} else {
			r.Links[SelfResource] = convertEchoPathToUriTemplatePath(path)
		}
		r.Meta.DocsUrl = fmt.Sprintf("%s/%s", r.Meta.DocsUrl, string(rc.ResourceName))
		r.Meta.Resource = rc.ResourceName
		r.StatusCode = httpStatus
		r.Success = success
		if si, ok := js.(SuccessInspector); ok {
			r.Success = si.IsSuccess()
		}
		var j []byte
		j, err = json.MarshalIndent(r, "", "   ")
		if err != nil {
			break
		}
		err = ctx.String(httpStatus, string(j))
		status = &Status{StatusCode: httpStatus}
	}
	if status == nil && err != nil {
		status = &Status{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}
	if status.Error != nil {
		b, _ := json.Marshal(status.ToResponse())
		_ = ctx.String(status.StatusCode, string(b))
	}
	return status
}
