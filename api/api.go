package api

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/config"
	"gearbox/global"
	"gearbox/jsonapi" // @TODO Refactor this out to interface{}s
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/util"
	"github.com/gedex/inflector"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

const Rootname types.RouteName = "root"

var _ Apier = (*Api)(nil)

type Api struct {
	Config        config.Configer
	Port          string
	Echo          *echo.Echo
	Parent        interface{}
	ControllerMap apimodeler.ControllerMap
}

func (me *Api) GetRootLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap {
	lm := make(apimodeler.LinkMap, 0)
	for k, ms := range me.ControllerMap {
		if k == apimodeler.Basepath {
			continue
		}
		s := ms
		lm.AddLink(
			GetQualifiedRelType(apimodeler.RelType(s.GetName())),
			apimodeler.Link(s.GetBasepath()),
		)
	}
	return lm
}

func (me *Api) GetListLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap {
	lm := make(apimodeler.LinkMap, 0)
	for range only.Once {
		path := me.GetSelfPath(ctx)
		if types.Basepath(path) == apimodeler.Basepath {
			break
		}
		if !ctx.Controller.CanAddItem(ctx) {
			break
		}
		lm.AddLink(
			GetQualifiedRelType(apimodeler.AddItemRelType),
			apimodeler.Link(fmt.Sprintf("%s/new", path)),
		)
		lm.AddLink(
			GetQualifiedRelType(apimodeler.ListRelType),
			apimodeler.Link(fmt.Sprintf("%s", path)),
		)
	}
	return lm
}

func (me *Api) GetItemLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap {
	lm := make(apimodeler.LinkMap, 0)
	lm.AddLink(
		GetQualifiedRelType(apimodeler.ItemRelType),
		apimodeler.Link(fmt.Sprintf("%s", me.GetSelfPath(ctx))),
	)
	return lm
}

func (me *Api) GetCommonLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap {
	lm := make(apimodeler.LinkMap, 0)
	lm.AddLink(
		GetQualifiedRelType(apimodeler.RootRelType),
		apimodeler.Link(apimodeler.Basepath),
	)
	lm.AddLink(
		GetQualifiedRelType(apimodeler.CurrentRelType),
		apimodeler.Link(me.GetSelfPath(ctx)),
	)
	return lm
}

type ConfigGetter interface {
	GetConfig() config.Configer
}

func NewApi(parent interface{}) *Api {
	c, ok := parent.(ConfigGetter)
	if !ok {
		panic("parent does not implement ConfigGetter")
	}

	a := Api{
		Config:        c.GetConfig(),
		Echo:          newConfiguredEcho(),
		Parent:        parent,
		ControllerMap: make(apimodeler.ControllerMap, 0),
	}
	a.Port = Port
	return &a
}

func newConfiguredEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// @TODO short circuit error handler, if needed
		var ok = false
		if !ok {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
	return e
}

func (me *Api) WireRoutes() {
	for _, _c := range me.ControllerMap {

		// Copy to allow different values in closures
		ctlr := _c

		e := me.Echo

		prefix := string(apimodeler.GetRouteNamePrefix(ctlr))

		var route *echo.Route

		path := string(apimodeler.GetBasepath(ctlr))

		// Collection Route
		route = e.GET(path, func(ctx echo.Context) (err error) {
			for range only.Once {
				rd := ja.NewRootDocument(ctx, global.ListResponse)
				c := apimodeler.NewContext(&apimodeler.ContextArgs{
					Contexter:      ctx,
					RootDocumentor: rd,
					Controller:     ctlr,
				})
				data, sts := ctlr.GetList(c)
				if is.Error(sts) {
					break
				}
				sts = me.setListData(c, data)
				if is.Error(sts) {
					break
				}
				lm, sts := ctlr.GetListLinkMap(c)
				if is.Error(sts) {
					break
				}
				for rt, lnk := range lm {
					rd.AddLink(rt, lnk)
				}
				err = me.JsonMarshalHandler(c, sts)
			}
			return err
		})
		if path == string(apimodeler.NoFilterPath) {
			route.Name = string(Rootname)
		} else {
			route.Name = fmt.Sprintf("%s-list", inflector.Pluralize(prefix))
		}

		urlTemplate := string(apimodeler.GetResourceUrlTemplate(ctlr))
		if urlTemplate == string(apimodeler.Basepath) {
			continue
		}

		// Single Item Route
		route = e.GET(urlTemplate, func(ec echo.Context) error {
			var sts status.Status
			rd := ja.NewRootDocument(ec, global.ItemResponse)
			ctx := apimodeler.NewContext(&apimodeler.ContextArgs{
				Contexter:      ec,
				RootDocumentor: rd,
				Controller:     ctlr,
			})
			for range only.Once {
				id, sts := apimodeler.GetIdFromUrl(ctx, ctlr)
				if is.Error(sts) {
					break
				}
				var ro *ja.ResourceObject
				ro, sts = getResourceObject(rd)
				if is.Error(sts) {
					break
				}
				var item apimodeler.ItemModeler
				item, sts = ctlr.GetItemDetails(ctx, id)
				if is.Error(sts) {
					break
				}
				var list apimodeler.List
				list, sts = item.GetRelatedItems(ctx)
				if is.Error(sts) {
					break
				}
				sts = me.setItemData(ctx, ro, item, list)
				rd.Data = ro
			}
			return me.JsonMarshalHandler(ctx, sts)
		})
		route.Name = fmt.Sprintf("%s-details", inflector.Singularize(prefix))
	}
}

func (me *Api) GetItemUrl(ctx *apimodeler.Context, item apimodeler.ItemModeler) (u types.UrlTemplate, sts status.Status) {
	for range only.Once {
		//
		// @TODO This may need to be make more robust later
		//
		u = types.UrlTemplate(fmt.Sprintf("%s/%s",
			// me.GetBaseUrl(),
			ctx.Controller.GetBasepath(),
			item.GetId(),
		))
	}
	return u, sts
}

func (me *Api) GetSelfUrl(ctx *apimodeler.Context) types.UrlTemplate {
	r := ctx.Request()
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, r.Host, me.GetSelfPath(ctx))
	return types.UrlTemplate(url)
}

func (me *Api) GetSelfPath(ctx *apimodeler.Context) types.UrlTemplate {
	return types.UrlTemplate(ctx.Request().RequestURI)
}

func (me *Api) GetContentType(ctx *apimodeler.Context) apimodeler.HttpHeaderValue {
	return ja.ContentType + "; " + apimodeler.CharsetUTF8
}

const (
	GearboxApiIdentifier         apimodeler.Metaname = "GearboxAPI"
	GearboxApiSchema             apimodeler.Link     = "https://docs.gearbox.works/api/schema/1.0/"
	MetaGearboxBaseurl           apimodeler.Metaname = GearboxApiIdentifier + ".baseurl"
	MetaGearboxApiSchema         apimodeler.Metaname = GearboxApiIdentifier + ".schema"
	SchemaGearboxApiRelationType                     = apimodeler.RelType("schema." + GearboxApiIdentifier)
)

func (me *Api) JsonMarshalHandler(ctx *apimodeler.Context, sts status.Status) status.Status {
	for range only.Once {
		sts = ctx.SetResponseHeader(echo.HeaderContentType, me.GetContentType(ctx))
		if is.Error(sts) {
			break
		}
		_, ok := ctx.GetRootDocument().(*ja.RootDocument)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "context.RootDocument() does not implement ja.RootDocument",
			})
			break
		}
		route, sts := me.getRouteName(ctx)
		if is.Error(sts) {
			break
		}
		if route == Rootname {
			ctx.AddLinks(me.GetRootLinkMap(ctx))
		}
		ctx.AddMeta(MetaGearboxApiSchema, route)
	}

	switch ctx.GetResponseType() {
	case global.ListResponse:
		ctx.AddLinks(me.GetListLinkMap(ctx))
	case global.ItemResponse:
		ctx.AddLinks(me.GetItemLinkMap(ctx))
	}

	ctx.AddLinks(me.GetCommonLinkMap(ctx))

	ctx.AddLink(apimodeler.SelfRelType, apimodeler.Link(me.GetSelfPath(ctx)))
	ctx.AddLink(apimodeler.SchemaDcRelType, apimodeler.DcSchema)
	ctx.AddLink(apimodeler.SchemaDcTermsRelType, apimodeler.DcTermsSchema)
	ctx.AddLink(SchemaGearboxApiRelationType, GearboxApiSchema)

	ctx.AddMeta(apimodeler.MetaDcCreator, GearboxApiIdentifier)
	ctx.AddMeta(apimodeler.MetaDcTermsIdentifier, me.GetSelfUrl(ctx))
	ctx.AddMeta(apimodeler.MetaDcLanguage, apimodeler.DefaultLanguage)
	ctx.AddMeta(MetaGearboxBaseurl, me.GetBaseUrl())

	if is.Error(sts) {
		_ = ctx.SetResponseStatus(sts.HttpStatus())
		ctx.SetErrors(sts)
	}
	sts = ctx.SendResponse()
	return sts
}

func (me *Api) AddController(controller apimodeler.ListController) (sts status.Status) {
	for range only.Once {
		getter, ok := controller.(apimodeler.BasepathGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "model does not implement apimodeler.BasepathGetter",
			})
			break
		}
		path := getter.GetBasepath()
		me.ControllerMap[path] = controller
	}
	return sts
}

//func (me *Apier) GetMethodMap() api.MethodMap {
//	return me.MethodMap
//}

func (me *Api) SetParent(parent interface{}) {
	me.Parent = parent
}

func (me *Api) GetBaseUrl() (url types.UrlTemplate) {
	return types.UrlTemplate(fmt.Sprintf(BaseUrlPattern, me.Port))
}

func (me *Api) Start() {
	err := me.Echo.Start(":" + me.Port)
	if err != nil {
		util.Error(err)
	}
}

func (me *Api) Stop() {
	err := me.Echo.Close()
	if err != nil {
		util.Error(err)
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

func (me *Api) setItemData(ctx *apimodeler.Context, ro *ja.ResourceObject, item apimodeler.ItemModeler, list apimodeler.List) (sts status.Status) {
	for range only.Once {
		sts = ro.SetId(item.GetId())
		if is.Error(sts) {
			break
		}
		sts = ro.SetType(item.GetType())
		if is.Error(sts) {
			break
		}
		sts = ro.SetAttributes(item)
		if is.Error(sts) {
			break
		}

		sts = ro.SetRelatedItems(ctx, list)
		if is.Error(sts) {
			break
		}

		if ctx.GetResponseType() == global.ItemResponse {
			break
		}
		var su types.UrlTemplate
		su, sts = me.GetItemUrl(ctx, item)
		if is.Error(sts) {
			break
		}
		getter, ok := item.(apimodeler.ItemLinkMapGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "item does not implement apimodeler.ItemLinkMapGetter",
			})
			break
		}
		lm, sts := getter.GetItemLinkMap(ctx)
		if is.Error(sts) {
			break
		}
		lm.AddLink(apimodeler.SelfRelType, apimodeler.Link(su))
		ro.SetLinks(lm)
	}
	return sts
}

func (me *Api) setListData(ctx *apimodeler.Context, data interface{}) (sts status.Status) {
	for range only.Once {
		getter, ok := data.(apimodeler.ListGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "data does not implement apimodeler.ListGetter",
			})
			break
		}
		var coll apimodeler.List
		coll, sts = getter.GetList(ctx)
		if is.Error(sts) {
			break
		}
		for _, item := range coll {
			ro := ja.NewResourceObject()
			sts = me.setItemData(ctx, ro, item, nil)
			if is.Error(sts) {
				break
			}
			sts = ctx.AddResponseItem(ro)
			if is.Error(sts) {
				break
			}
		}
	}
	return sts
}

func (me *Api) getRouteName(ctx *apimodeler.Context) (name types.RouteName, sts status.Status) {
	for range only.Once {
		rts := me.Echo.Routes()
		path, sts := ctx.GetRequestPath()
		if is.Error(sts) {
			break
		}
		for _, rt := range rts {
			if rt.Path == string(path) {
				name = types.RouteName(rt.Name)
				break
			}
		}
	}
	return name, sts
}

func GetQualifiedRelType(reltype apimodeler.RelType) apimodeler.RelType {
	rt := fmt.Sprintf(apimodeler.RelTypePattern, GearboxApiIdentifier, reltype)
	return apimodeler.RelType(rt)
}
