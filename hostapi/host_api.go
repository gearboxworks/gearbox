package hostapi

import (
	"fmt"
	"gearbox"
	"gearbox/api"
	"gearbox/apibuilder"
	"gearbox/config"
	"gearbox/jsonapi"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"github.com/gedex/inflector"
	"github.com/labstack/echo"
)

var _ gearbox.HostApi = (*HostApi)(nil)

type HostApi struct {
	Config         config.Configer
	Api            *api.Api
	Gearbox        gearbox.Gearboxer
	ConnectionsMap ab.ConnectionsMap
}

func apiResponseDefaults() *api.Response {
	return &api.Response{
		Meta: api.ResponseMeta{
			Service: ServiceName,
			Version: Version,
			DocsUrl: DocsUrl,
		},
		Links: make(api.Links, 0),
	}
}

func NewHostApi(gearbox gearbox.Gearboxer) *HostApi {
	ha := &HostApi{
		Config:         gearbox.GetConfig(),
		Api:            api.NewApi(echo.New(), apiResponseDefaults()),
		Gearbox:        gearbox,
		ConnectionsMap: make(ab.ConnectionsMap, 0),
	}
	ha.Api.Port = Port
	return ha
}

func (me *HostApi) Route() (sts status.Status) {
	for range only.Once {
		sts = me.addRoutes()
		if is.Error(sts) {
			return sts
		}
		me.connectRoutes(me.ConnectionsMap)
	}
	return sts
}

func (me *HostApi) connectRoutes(connectionsMap ab.ConnectionsMap) {
	for _, cn := range connectionsMap {
		e := me.Api.Echo

		prefix := cn.GetRouteNamePrefix()

		var route *echo.Route

		// Collection Route
		route = e.GET(string(cn.GetBasepath()), func(ctx echo.Context) error {
			rd := ja.NewRootDocument(ja.CollectionResponse)
			data, sts := ab.GetCollectionSlice(cn.Self.GetCollection(ab.NoFilterPath))
			if is.Success(sts) {
				sts = setCollectionData(rd, data)
			}
			return me.JsonMarshalHandler(rd, ctx, sts)
		})
		route.Name = fmt.Sprintf("%s-list", inflector.Pluralize(prefix))

		// Single Item Route
		route = e.GET(string(cn.GetResourceUrlTemplate()), func(ctx echo.Context) error {
			var sts status.Status
			rd := ja.NewRootDocument(ja.DatasetResponse)
			for range only.Once {
				id, sts := cn.GetIdFromUrl(ctx)
				if is.Error(sts) {
					break
				}
				var ro *ja.ResourceObject
				ro, sts = getResourceObject(rd)
				if is.Error(sts) {
					break
				}
				var item ab.Item
				item, sts = cn.Self.GetItem(id, ctx)
				if is.Error(sts) {
					break
				}
				sts = setItemData(ro, item)
				rd.Data = ro
			}
			return me.JsonMarshalHandler(rd, ctx, sts)
		})
		route.Name = fmt.Sprintf("%s-details", inflector.Singularize(prefix))

	}
}

func getResourceObject(rd *ja.RootDocument) (ro *ja.ResourceObject, sts status.Status) {
	ro, ok := rd.Data.(*ja.ResourceObject)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: "root document does not contain a single resource object",
		})
	}
	return ro, sts
}

func setItemData(ro *ja.ResourceObject, item ab.Item) (sts status.Status) {
	for range only.Once {
		itemId := item.GetId()
		sts = ro.SetId(ja.ResourceId(itemId))
		if is.Error(sts) {
			break
		}
		typ := item.GetType()
		sts = ro.SetType(ja.ResourceType(typ))
		if is.Error(sts) {
			break
		}
		sts = ro.SetAttributes(item)
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func setCollectionData(rd *ja.RootDocument, data interface{}) (sts status.Status) {
	for range only.Once {
		getter, ok := data.(ab.ItemsGetter)
		if !ok {
			break
		}
		var items ab.Collection
		items, sts = getter.GetItems()
		if is.Error(sts) {
			break
		}
		for _, item := range items {
			ro := ja.NewResourceObject()
			sts = setItemData(ro, item)
			if is.Error(sts) {
				break
			}
			sts = rd.AddResourceObject(ro)
			if is.Error(sts) {
				break
			}
		}
	}
	return sts
}

func (me *HostApi) getRouteName(ctx echo.Context) (name api.RouteName) {
	rts := me.Api.Echo.Routes()
	path := ctx.Path()
	for _, rt := range rts {
		if rt.Path == path {
			name = api.RouteName(rt.Name)
			break
		}
	}
	return name
}

func (me *HostApi) GetSelfUrl(ctx echo.Context) ab.UrlTemplate {
	r := ctx.Request()
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, r.Host, me.GetSelfPath(ctx))
	return ab.UrlTemplate(url)
}

func (me *HostApi) GetSelfPath(ctx echo.Context) ab.UrlTemplate {
	return ab.UrlTemplate(ctx.Request().RequestURI)
}

func (me *HostApi) GetContentType(ctx echo.Context) string {
	return ja.ContentType + "; " + ja.CharsetUTF8
}

const (
	GearboxApiSchema             ja.Link         = "https://docs.gearbox.works/api/schema/1.0/"
	SchemaGearboxApiRelationType ja.RelationType = "schema.GearboxAPI"
	MetaGearboxApiResponseType   ja.Metaname     = "GearboxAPI.response_type"
)

func (me *HostApi) JsonMarshalHandler(rootdoc *ja.RootDocument, ctx echo.Context, sts status.Status) status.Status {
	for range only.Once {
		var err error
		ctx.Response().Header().Set(echo.HeaderContentType, me.GetContentType(ctx))

		rootdoc.MetaMap[ja.MetaDcCreator] = gearbox.Brandname

		rootdoc.MetaMap[ja.MetaDcTermsIdentifier] = me.GetSelfUrl(ctx)
		rootdoc.MetaMap[ja.MetaDcLanguage] = ja.DefaultLanguage

		rootdoc.LinkMap[ja.SelfRelationType] = ja.Link(me.GetSelfPath(ctx))
		rootdoc.LinkMap[ja.SchemaDcRelationType] = ja.DcSchema
		rootdoc.LinkMap[ja.SchemaDcTermsRelationType] = ja.DcTermsSchema

		rootdoc.LinkMap[SchemaGearboxApiRelationType] = GearboxApiSchema
		rootdoc.MetaMap[MetaGearboxApiResponseType] = me.getRouteName(ctx)

		if is.Error(sts) {
			ctx.Response().Status = sts.HttpStatus()
			_ = rootdoc.SetError(sts)
		}
		err = ctx.JSONPretty(ctx.Response().Status, rootdoc, "   ")
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("error when sending output for '%s'",
					ctx.Path(),
				),
			})
			break
		}
	}
	return sts
}

func (me *HostApi) AddConnector(connector ab.Connector) (sts status.Status) {
	for range only.Once {
		getter, ok := connector.(ab.BasepathGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "factory has no GetBasepath()",
			})
			break
		}
		me.ConnectionsMap[getter.GetBasepath()] = ab.NewConnections(connector)
	}
	return sts
}

//func (me *HostApi) GetValuesFunc(name api.RouteName) (values api.ValuesFunc, sts status.Status) {
//	for range only.Once {
//		var ok bool
//		values, ok = me.Api.ValuesFuncMap[name]
//		if !ok {
//			sts = status.Fail(&status.Args{
//				Message: "no values func for route '%s'",
//			})
//		}
//	}
//	return values, sts
//
//}
//
//func (me *HostApi) GetUriTemplateVars(name api.RouteName, values interface{}, index int) (api.UriTemplateVars, status.Status) {
//	return me.Api.GetUriTemplateVars(name, values, index)
//}

func (me *HostApi) GetMethodMap() api.MethodMap {
	return me.Api.MethodMap
}

func (me *HostApi) SetGearbox(gb gearbox.Gearboxer) {
	me.Gearbox = gb
}

func (me *HostApi) GetBaseUrl() (url api.UriTemplate) {
	return me.Api.GetBaseUrl()
}

func (me *HostApi) GetUrl(name api.RouteName, vars api.UriTemplateVars) (url api.UriTemplate, sts status.Status) {
	return me.Api.GetUrl(name, vars)
}

func (me *HostApi) GetUrlPath(name api.RouteName, vars api.UriTemplateVars) (url api.UriTemplate, sts status.Status) {
	return me.Api.GetUrlPath(name, vars)
}

func (me *HostApi) GetUrlPathTemplate(name api.RouteName) (url api.UriTemplate, sts status.Status) {
	for range only.Once {
		if me.Api == nil {
			sts = status.NewStatus(&status.Args{
				Message: fmt.Sprintf("accessing host api when internal api property is nil for resource type '%s'",
					name,
				),
				Help: status.ContactSupportHelp(),
			})
			break
		}
		url, sts = me.Api.GetUrlPathTemplate(name)
		if status.IsError(sts) {
			break
		}
	}
	return url, sts
}

func (me *HostApi) Url() string {
	return fmt.Sprintf("http://127.0.0.1:%s", me.Api.Port)
}

func (me *HostApi) Start() {
	me.Api.Start()
}

func (me *HostApi) Relate(primary api.RouteName, related *api.Related) {
	me.Api.Relate(primary, related)
}

func (me *HostApi) Stop() {
	me.Api.Stop()
}

func (me *HostApi) GET___(path api.UriTemplate, name api.RouteName, funcs ...interface{}) *echo.Route {
	handler, valuesFunc := getFuncsArgs(funcs)
	return me.Api.GET(path, name, valuesFunc, func(rc *api.RequestContext) (err error) {
		me.Gearbox.SetRouteName(rc.RouteName)
		if handler != nil {
			err = rc.JsonMarshalHandler(handler(rc))
		} else {
			err = rc.JsonMarshalHandler(nil)
		}
		return err
	})
}

func (me *HostApi) POST__(path api.UriTemplate, name api.RouteName, funcs ...interface{}) *echo.Route {
	handler, valuesFunc := getFuncsArgs(funcs)
	return me.Api.POST(path, name, valuesFunc, func(rc *api.RequestContext) error {
		return rc.JsonMarshalHandler(handler(rc))
	})
}

func (me *HostApi) PUT___(path api.UriTemplate, name api.RouteName, funcs ...interface{}) *echo.Route {
	handler, valuesFunc := getFuncsArgs(funcs)
	return me.Api.PUT(path, name, valuesFunc, func(rc *api.RequestContext) error {
		return rc.JsonMarshalHandler(handler(rc))
	})
}

func (me *HostApi) DELETE(path api.UriTemplate, name api.RouteName, funcs ...interface{}) *echo.Route {
	handler, valuesFunc := getFuncsArgs(funcs)
	return me.Api.DELETE(path, name, valuesFunc, func(rc *api.RequestContext) error {
		return rc.JsonMarshalHandler(handler(rc))
	})
}

func getFuncsArgs(funcs []interface{}) (handler api.UpstreamHandlerFunc, valuesFunc api.ValuesFunc) {
	switch {
	case len(funcs) > 1:
		if funcs[1] != nil {
			valuesFunc = funcs[1].(func(...interface{}) (api.ValuesFuncValues, status.Status))
		}
		fallthrough
	case len(funcs) > 0:
		if funcs[0] != nil {
			handler = funcs[0].(func(*api.RequestContext) interface{})
		}
	}
	return handler, valuesFunc
}
