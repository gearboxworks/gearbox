package api

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestContext struct {
	*Api
	RouteName
	echo.Context
}

func NewRequestContext(_api *Api, routeName RouteName) *RequestContext {
	return &RequestContext{
		Api:       _api,
		RouteName: routeName,
	}
}
func (me *RequestContext) Param(name ResourceVarName) (value string) {
	if me.Context != nil {
		value = me.Context.Param(string(name))
	}
	return value
}

func (me *RequestContext) WrapHandler(next HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		me.Context = ctx
		return next(me)
	}
}

func (me *RequestContext) ReadRequestBody() (b []byte, sts status.Status) {
	b, err := ioutil.ReadAll(me.Context.Request().Body)
	if err != nil {
		sts = status.Wrap(err, &status.Args{
			Message: fmt.Sprintf("failed to read request body for '%s'", me.RouteName),
			Help:    ContactSupportHelp(),
		})
	}
	return b, sts
}

func (me *RequestContext) CloseRequestBody() {
	_ = me.Context.Request().Body.Close()
}

func (me *RequestContext) GetApiSelfLink() (path UriTemplate) {
	path = UriTemplate(me.Context.Path())
	parts := strings.Split(string(path), "/")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		if p[0] != ':' {
			continue
		}
		p = p[1:]
		parts[i] = me.Param(ResourceVarName(p))
	}
	path = UriTemplate(strings.Join(parts, "/"))
	return path
}

func (me *RequestContext) UnmarshalFromRequest(obj interface{}) (sts status.Status) {
	for range only.Once {
		b, sts := me.ReadRequestBody()
		if status.IsError(sts) {
			break
		}
		err := json.Unmarshal(b, &obj)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message:    fmt.Sprintf("unexpected format for request body: '%s'", string(b)),
				HttpStatus: http.StatusBadRequest,
				Help:       GetApiHelp(me.RouteName),
			})
			break
		}
	}
	me.CloseRequestBody()
	return sts
}

func (me *RequestContext) getResourceLinks(data interface{}, inlinks Links) (outlinks Links, sts status.Status) {
	for range only.Once {
		outlinks = make(Links, 0)
		for rn, lnk := range inlinks {
			var ut UriTemplate
			s, ok := lnk.(string)
			if ok {
				ut = UriTemplate(s)
			} else {
				ut, ok = lnk.(UriTemplate)
			}
			if !ok {
				continue
			}
			if !strings.HasPrefix(string(ut), me.Context.Path()) {
				continue
			}
			if s == me.Context.Path() {
				continue
			}
			vf, ok := me.Api.ValuesFuncMap[rn]
			if !ok {
				continue
			}
			if vf == nil {
				continue
			}
			values, sts := vf()
			if status.IsError(sts) {
				sts = status.Wrap(sts, &status.Args{
					Message: fmt.Sprintf("unable to acccess links for resource '%s'", me.RouteName),
				})
				break
			}
			if len(values) == 0 {
				continue
			}
			for i := range values[0] {
				var vars UriTemplateVars
				vars, sts = me.Api.GetUriTemplateVars(rn, values, i)
				if status.IsError(sts) {
					// @TODO log the error
					continue
				}
				if len(vars) == 0 {
					continue
				}
				var url UriTemplate
				url, sts = me.GetUrlPath(rn, vars)
				if status.IsError(sts) {
					// @TODO log the error
					continue
				}
				ln := RouteName(fmt.Sprintf("%s%s", rn, vars.String()))
				outlinks[ln] = url
			}
		}
	}
	return outlinks, sts
}

func (me *RequestContext) JsonMarshalHandler(data interface{}) (sts status.Status) {
	var err error
	ctx := me.Context
	ctx.Set(RequestContextKey, me) // Used in NewApi to short-circuit default error handling
	for range only.Once {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		var ok bool
		success := true
		sts, ok = data.(status.Status)
		if ok {
			success = false
			ctx.Response().Status = sts.HttpStatus()
		}
		httpStatus := ctx.Response().Status
		r := *me.Api.Defaults.Clone()
		if rdg, ok := data.(ResponseDataGetter); !ok {
			r.Data = data
		} else {
			r.Data = rdg.(ResponseDataGetter).GetResponseData()
		}
		path := ctx.Path()
		if path != "/" {
			var links Links
			links, sts = me.getResourceLinks(data, r.Links)
			r.Links = make(Links, 0)
			mmg := me.Api.MethodMap[Method(http.MethodGet)]
			r.Links[LinksResource] = "/"
			r.Links[MetaEndpointsResource] = mmg[MetaEndpointsResource]
			r.Links[MetaMethodsResource] = mmg[MetaMethodsResource]
			if len(links) > 0 && status.IsSuccess(sts) {
				r.Links[me.RouteName] = links
			}
		}
		if slg, ok := data.(UrlPathGetter); ok {
			url, _ := slg.GetApiUrlPath()
			r.Links[SelfResource] = url
		} else {
			r.Links[SelfResource] = me.GetApiSelfLink()
		}
		r.Meta.DocsUrl = fmt.Sprintf("%s/%s", r.Meta.DocsUrl, string(me.RouteName))
		r.Meta.RouteName = me.RouteName
		r.StatusCode = httpStatus
		r.Success = success
		if si, ok := data.(status.SuccessInspector); ok {
			r.Success = si.IsSuccess()
		}
		var j []byte
		// @TODO Add ?format=yes to pretty print JSON
		j, err = json.MarshalIndent(r, "", "   ")
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marshal output for resource '%s'",
					me.RouteName,
				),
			})
			break
		}
		err = ctx.String(httpStatus, string(j))
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("error when sending output for '%s'",
					me.RouteName,
				),
			})
			break
		}
	}
	return sts
}
