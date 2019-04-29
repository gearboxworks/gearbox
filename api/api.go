package api

import (
	"fmt"
	"gearbox/apiworks"
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
	"io/ioutil"
	"log"
	"net/http"
)

const Rootname types.RouteName = "root"

var _ Apier = (*Api)(nil)

type Api struct {
	Config        config.Configer
	Port          string
	Echo          *echo.Echo
	Parent        interface{}
	ControllerMap apiworks.ControllerMap
}

func (me *Api) GetRootLinkMap(ctx *apiworks.Context) apiworks.LinkMap {
	lm := make(apiworks.LinkMap, 0)
	for k, ms := range me.ControllerMap {
		if k == apiworks.Basepath {
			continue
		}
		s := ms
		lm.AddLink(
			GetQualifiedRelType(apiworks.RelType(s.GetName())),
			apiworks.Link(s.GetBasepath()),
		)
	}
	return lm
}

func (me *Api) GetListLinkMap(ctx *apiworks.Context) apiworks.LinkMap {
	lm := make(apiworks.LinkMap, 0)
	for range only.Once {
		path := me.GetSelfPath(ctx)
		if types.Basepath(path) == apiworks.Basepath {
			break
		}
		if !ctx.Controller.CanAddItem(ctx) {
			break
		}
		lm.AddLink(
			GetQualifiedRelType(apiworks.AddItemRelType),
			apiworks.Link(fmt.Sprintf("%s/new", path)),
		)
		lm.AddLink(
			GetQualifiedRelType(apiworks.ListRelType),
			apiworks.Link(fmt.Sprintf("%s", path)),
		)
	}
	return lm
}

func (me *Api) GetItemLinkMap(ctx *apiworks.Context) apiworks.LinkMap {
	lm := make(apiworks.LinkMap, 0)
	lm.AddLink(
		GetQualifiedRelType(apiworks.ItemRelType),
		apiworks.Link(fmt.Sprintf("%s", me.GetSelfPath(ctx))),
	)
	return lm
}

func (me *Api) GetCommonLinkMap(ctx *apiworks.Context) apiworks.LinkMap {
	lm := make(apiworks.LinkMap, 0)
	lm.AddLink(
		GetQualifiedRelType(apiworks.RootRelType),
		apiworks.Link(apiworks.Basepath),
	)
	lm.AddLink(
		GetQualifiedRelType(apiworks.CurrentRelType),
		apiworks.Link(me.GetSelfPath(ctx)),
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
		ControllerMap: make(apiworks.ControllerMap, 0),
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

		prefix := string(apiworks.GetRouteNamePrefix(ctlr))

		listpath := string(apiworks.GetBasepath(ctlr))
		me.WireListRoute(
			me.Echo,
			ctlr,
			prefix,
			listpath,
		)

		itempath := string(apiworks.GetResourceUrlTemplate(ctlr))
		if itempath == string(apiworks.Basepath) {
			continue
		}

		me.WireNewItemRoute(
			me.Echo,
			ctlr,
			prefix,
			fmt.Sprintf("%s/new", listpath),
		)

		me.WireItemRoute(
			me.Echo,
			ctlr,
			prefix,
			itempath,
		)

		me.WireItemDeleteRoute(
			me.Echo,
			ctlr,
			prefix,
			itempath,
		)

	}
}

func (me *Api) WireListRoute(e *echo.Echo, lc apiworks.ListController, prefix, path string) {

	route := e.GET(path, func(ec echo.Context) (err error) {
		for range only.Once {
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ListResponse)
			data, sts := lc.GetList(ctx)
			if is.Error(sts) {
				break
			}
			sts = me.setListData(ctx, data)
			if is.Error(sts) {
				break
			}
			lm, sts := lc.GetListLinkMap(ctx)
			if is.Error(sts) {
				break
			}
			for rt, lnk := range lm {
				rd.AddLink(rt, lnk)
			}
			err = me.JsonMarshalHandler(ctx, sts)
		}
		return err
	})
	if path == string(apiworks.NoFilterPath) {
		route.Name = string(Rootname)
	} else {
		route.Name = fmt.Sprintf("%s-list", inflector.Pluralize(prefix))
	}
}

func (me *Api) WireItemRoute(e *echo.Echo, lc apiworks.ListController, prefix, path string) {

	// Single Item Route
	route := e.GET(path, func(ec echo.Context) error {
		var sts status.Status
		rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
		for range only.Once {
			id, sts := apiworks.GetIdFromUrl(ctx, lc)
			if is.Error(sts) {
				break
			}
			var ro *jsonapi.ResourceObject
			ro, sts = getResourceObject(rd)
			if is.Error(sts) {
				break
			}
			var item apiworks.ItemModeler
			item, sts = lc.GetItemDetails(ctx, id)
			if is.Error(sts) {
				break
			}

			var list apiworks.List
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

func closeRequestBody(ec echo.Context) {
	err := ec.Request().Body.Close()
	if err != nil {
		log.Printf(
			"Could not close response body from HttpRequest: %s\n",
			err.Error(),
		)
	}
}

func (me *Api) WireNewItemRoute(e *echo.Echo, lc apiworks.ListController, prefix, path string) {

	// Single Item Route
	route := e.POST(path, func(ec echo.Context) error {
		var sts status.Status
		rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
		for range only.Once {
			defer closeRequestBody(ec)
			b, err := ioutil.ReadAll(ec.Request().Body)
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("unable to read body of '%s' request", path),
				})
				break
			}
			var ro jsonapi.ResourceObject
			err = ro.Unmarshal(b)
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message:    fmt.Sprintf("unable to unmarshal body of '%s' request", path),
					HttpStatus: http.StatusBadRequest,
				})
				break
			}
			sts = ctx.Controller.AddItem(ctx, &ro)
			if is.Error(sts) {
				break
			}
			rd.Data = ro.ResourceIdObject

		}
		return me.JsonMarshalHandler(ctx, sts)
	})
	route.Name = fmt.Sprintf("add-%s", inflector.Singularize(prefix))

}

func (me *Api) WireItemDeleteRoute(e *echo.Echo, lc apiworks.ListController, prefix, path string) {

	route := e.DELETE(path, func(ec echo.Context) error {
		var sts status.Status
		rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
		for range only.Once {
			var id apiworks.ItemId
			id, sts = apiworks.GetIdFromUrl(ctx, lc)
			if is.Error(sts) {
				break
			}
			sts = ctx.Controller.DeleteItem(ctx, id)
			if is.Error(sts) {
				break
			}
			rd.Data = jsonapi.NewResourceIdObjectWithIdType(
				jsonapi.ResourceId(id),
				jsonapi.ResourceType(
					lc.GetNilItem(ctx).GetType(),
				),
			)
		}
		return me.JsonMarshalHandler(ctx, sts)
	})
	route.Name = fmt.Sprintf("delete-%s", inflector.Singularize(prefix))

}

func getRootDocumentAndContext(ec echo.Context, lc apiworks.ListController, rt types.ResponseType) (rd *jsonapi.RootDocument, ctx *apiworks.Context) {
	rd = jsonapi.NewRootDocument(ec, rt)
	ctx = apiworks.NewContext(&apiworks.ContextArgs{
		Contexter:      ec,
		RootDocumentor: rd,
		Controller:     lc,
	})
	return rd, ctx
}

func (me *Api) GetItemUrl(ctx *apiworks.Context, item apiworks.ItemModeler) (u types.UrlTemplate, sts status.Status) {
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

func (me *Api) GetSelfUrl(ctx *apiworks.Context) types.UrlTemplate {
	r := ctx.Request()
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, r.Host, me.GetSelfPath(ctx))
	return types.UrlTemplate(url)
}

func (me *Api) GetSelfPath(ctx *apiworks.Context) types.UrlTemplate {
	return types.UrlTemplate(ctx.Request().RequestURI)
}

func (me *Api) GetContentType(ctx *apiworks.Context) apiworks.HttpHeaderValue {
	return jsonapi.ContentType + "; " + apiworks.CharsetUTF8
}

const (
	GearboxApiIdentifier         apiworks.Metaname = "GearboxAPI"
	GearboxApiSchema             apiworks.Link     = "https://docs.gearbox.works/api/schema/1.0/"
	MetaGearboxBaseurl           apiworks.Metaname = GearboxApiIdentifier + ".baseurl"
	MetaGearboxApiSchema         apiworks.Metaname = GearboxApiIdentifier + ".schema"
	SchemaGearboxApiRelationType                   = apiworks.RelType("schema." + GearboxApiIdentifier)
)

func (me *Api) JsonMarshalHandler(ctx *apiworks.Context, sts status.Status) status.Status {
	var _sts status.Status
	for range only.Once {
		_, ok := ctx.GetRootDocument().(*jsonapi.RootDocument)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "context.RootDocument() does not implement ja.RootDocument",
			})
			break
		}
		if is.Error(sts) {
			// If an error occurred during context generation
			// and prior to this func being called
			break
		}
		_sts = ctx.SetResponseHeader(echo.HeaderContentType, me.GetContentType(ctx))
		if is.Error(_sts) {
			break
		}
		route, _sts := me.getRouteName(ctx)
		if is.Error(_sts) {
			break
		}
		if route == Rootname {
			ctx.AddLinks(me.GetRootLinkMap(ctx))
		}
		ctx.AddMeta(MetaGearboxApiSchema, route)
	}
	if _sts != nil {
		sts = _sts
	}
	if is.Error(sts) {
		_ = ctx.SetResponseStatus(sts.HttpStatus())
		ctx.SetErrors(sts)
	} else {
		switch ctx.GetResponseType() {
		case global.ListResponse:
			ctx.AddLinks(me.GetListLinkMap(ctx))
		case global.ItemResponse:
			ctx.AddLinks(me.GetItemLinkMap(ctx))
		}
	}

	ctx.AddLinks(me.GetCommonLinkMap(ctx))

	ctx.AddLink(apiworks.SelfRelType, apiworks.Link(me.GetSelfPath(ctx)))
	ctx.AddLink(apiworks.SchemaDcRelType, apiworks.DcSchema)
	ctx.AddLink(apiworks.SchemaDcTermsRelType, apiworks.DcTermsSchema)
	ctx.AddLink(SchemaGearboxApiRelationType, GearboxApiSchema)

	ctx.AddMeta(apiworks.MetaDcCreator, GearboxApiIdentifier)
	ctx.AddMeta(apiworks.MetaDcTermsIdentifier, me.GetSelfUrl(ctx))
	ctx.AddMeta(apiworks.MetaDcLanguage, apiworks.DefaultLanguage)
	ctx.AddMeta(MetaGearboxBaseurl, me.GetBaseUrl())

	if sts == nil {
		sts = status.Success("")
	}
	sts = ctx.SendResponse(sts)
	return sts
}

func (me *Api) AddController(controller apiworks.ListController) (sts status.Status) {
	for range only.Once {
		getter, ok := controller.(apiworks.BasepathGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "model does not implement apiworks.BasepathGetter",
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

func getResourceObject(rd *jsonapi.RootDocument) (ro *jsonapi.ResourceObject, sts status.Status) {
	ro, ok := rd.Data.(*jsonapi.ResourceObject)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: "root document does not contain a single resource object",
		})
	}
	ro.Renew()
	return ro, sts
}

func (me *Api) setItemData(ctx *apiworks.Context, ro *jsonapi.ResourceObject, item apiworks.ItemModeler, list apiworks.List) (sts status.Status) {
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

		if list == nil {
			break
		}

		sts = ro.SetRelatedItems(ctx, item, list)
		if is.Error(sts) {
			break
		}

		sts = me.SetRelationshipLinkMap(ctx, item, ro)
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
		getter, ok := item.(apiworks.ItemLinkMapGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "item does not implement apiworks.ItemLinkMapGetter",
			})
			break
		}
		lm, sts := getter.GetItemLinkMap(ctx)
		if is.Error(sts) {
			break
		}
		lm.AddLink(apiworks.SelfRelType, apiworks.Link(su))
		ro.SetLinks(lm)
	}
	return sts
}

func (me *Api) SetRelationshipLinkMap(ctx *apiworks.Context, item apiworks.ItemModeler, ro *jsonapi.ResourceObject) (sts status.Status) {
	for range only.Once {
		lm, sts := ctx.RootDocumentor.GetDataRelationshipsLinkMap()
		if is.Error(sts) {
			break
		}
		name := item.GetType()

		for rel, link := range lm {

			rel = apiworks.RelType(fmt.Sprintf("%s.%s", name, rel))

			ctx.RootDocumentor.AddLink(GetQualifiedRelType(rel), link)
		}
	}
	return sts
}

func (me *Api) setListData(ctx *apiworks.Context, data interface{}) (sts status.Status) {
	for range only.Once {
		getter, ok := data.(apiworks.ListGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "data does not implement apiworks.ListGetter",
			})
			break
		}
		var coll apiworks.List
		coll, sts = getter.GetList(ctx)
		if is.Error(sts) {
			break
		}
		for _, item := range coll {
			ro := jsonapi.NewResourceObject()
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

func (me *Api) getRouteName(ctx *apiworks.Context) (name types.RouteName, sts status.Status) {
	for range only.Once {
		rts := me.Echo.Routes()
		path, sts := ctx.GetRequestTemplatePath()
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

func GetQualifiedRelType(reltype apiworks.RelType) apiworks.RelType {
	rt := fmt.Sprintf(apiworks.RelTypePattern, GearboxApiIdentifier, reltype)
	return apiworks.RelType(rt)
}
