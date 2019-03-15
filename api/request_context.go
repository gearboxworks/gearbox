package api

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strings"
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

func (me *RequestContext) ReadRequestBody() (b []byte, status stat.Status) {
	b, err := ioutil.ReadAll(me.Context.Request().Body)
	if err != nil {
		status = stat.NewFailedStatus(&stat.Args{
			Error:   err,
			Message: fmt.Sprintf("failed to read request body for '%s'", me.ResourceName),
			Help:    ContactSupportHelp(),
		})
	}
	return b, status
}

func (me *RequestContext) CloseRequestBody() {
	_ = me.Context.Request().Body.Close()
}

func (me *RequestContext) GetApiSelfLink() (path string) {
	path = me.Context.Path()
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		if p[0] != ':' {
			continue
		}
		p = p[1:]
		parts[i] = me.Param(p)
	}
	path = strings.Join(parts, "/")
	return path
}

func (me *RequestContext) UnmarshalFromRequest(obj interface{}) (status stat.Status) {
	for range only.Once {
		b, status := me.ReadRequestBody()
		if status.IsError() {
			break
		}
		err := json.Unmarshal(b, &obj)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:      err,
				Message:    fmt.Sprintf("unexpected format for request body: '%s'", string(b)),
				HttpStatus: http.StatusBadRequest,
				Help:       GetApiHelp(me.ResourceName),
			})
			break
		}
	}
	me.CloseRequestBody()
	return status
}

// @TODO Add ?format=yes to pretty print JSON
func (me *RequestContext) JsonMarshalHandler(js interface{}) (status *Status) {
	var err error
	ctx := me.Context
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
		r := *me.Api.Defaults.Clone()
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
			r.Links[SelfResource] = me.GetApiSelfLink()
		}
		r.Meta.DocsUrl = fmt.Sprintf("%s/%s", r.Meta.DocsUrl, string(me.ResourceName))
		r.Meta.Resource = me.ResourceName
		r.StatusCode = httpStatus
		r.Success = success
		if si, ok := js.(stat.SuccessInspector); ok {
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
