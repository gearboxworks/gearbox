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
	ControllerMap ControllerMap
}

func (me *Api) GetRootLinkMap(ctx *Context) LinkMap {
	lm := make(LinkMap, 0)
	for k, ms := range me.ControllerMap {
		if k == Basepath {
			continue
		}
		s := ms
		lm.AddLink(
			getQualifiedRelType(RelType(s.GetName())),
			Link(s.GetBasepath()),
		)
	}
	return lm
}

func (me *Api) GetListLinkMap(ctx *Context) LinkMap {
	lm := make(LinkMap, 0)
	for range only.Once {
		path := me.GetSelfPath(ctx)
		if types.Basepath(path) == Basepath {
			break
		}
		if !ctx.Controller.CanAddItem(ctx) {
			break
		}
		lm.AddLink(
			getQualifiedRelType(AddItemRelType),
			Link(fmt.Sprintf("%s/new", path)),
		)
		lm.AddLink(
			getQualifiedRelType(ListRelType),
			Link(fmt.Sprintf("%s", path)),
		)
	}
	return lm
}

func (me *Api) GetItemLinkMap(ctx *Context) LinkMap {
	lm := make(LinkMap, 0)
	lm.AddLink(
		getQualifiedRelType(ItemRelType),
		Link(fmt.Sprintf("%s", me.GetSelfPath(ctx))),
	)
	return lm
}

func (me *Api) GetCommonLinkMap(ctx *Context) LinkMap {
	lm := make(LinkMap, 0)
	lm.AddLink(
		getQualifiedRelType(RootRelType),
		Link(Basepath),
	)
	lm.AddLink(
		getQualifiedRelType(CurrentRelType),
		Link(me.GetSelfPath(ctx)),
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
		ControllerMap: make(ControllerMap, 0),
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

		route := me.WireListRoute(
			me.Echo,
			ctlr,
			listpath,
		)
		if listpath == string(NoFilterPath) {
			route.SetName(Rootname)
		} else {
			route.SetSingularName("get-%s-list", prefix)
		}

		itempath := string(apiworks.GetResourceUrlTemplate(ctlr))
		if itempath == string(Basepath) {
			continue
		}

		me.WireGetItemRoute(
			me.Echo,
			ctlr,
			itempath,
		).SetSingularName("get-%s", prefix)

		me.WireNewItemRoute(
			me.Echo,
			ctlr,
			prefix,
			fmt.Sprintf("%s/new", listpath),
		).SetSingularName("add-%s", prefix)

		me.WireUpdateItemRoute(
			me.Echo,
			ctlr,
			itempath,
		).SetSingularName("update-%s", prefix)

		me.WireDeleteItemRoute(
			me.Echo,
			ctlr,
			itempath,
		).SetSingularName("delete-%s", prefix)

	}
}

func (me *Api) WireListRoute(e *echo.Echo, lc ListController, path string) *Route {
	return &Route{
		Route: e.GET(path, func(ec echo.Context) (err error) {
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
		}),
	}
}

func (me *Api) WireGetItemRoute(e *echo.Echo, lc ListController, path string) *Route {

	return &Route{
		Route: e.GET(path, func(ec echo.Context) error {
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
				var item ItemModeler
				item, sts = lc.GetItemDetails(ctx, id)
				if is.Error(sts) {
					break
				}

				var list List
				list, sts = item.GetRelatedItems(ctx)
				if is.Error(sts) {
					break
				}
				sts = me.setItemData(ctx, ro, item, list)
				rd.Data = ro
			}
			return me.JsonMarshalHandler(ctx, sts)
		}),
	}
}

func closeRequestBody(ctx Contexter) {
	err := ctx.Request().Body.Close()
	if err != nil {
		log.Printf(
			"Could not close response body from HttpRequest: %s\n",
			err.Error(),
		)
	}
}
func (me *Api) WireNewItemRoute(e *echo.Echo, lc ListController, prefix, path string) *Route {
	return &Route{
		Route: e.POST(path, func(ec echo.Context) error {
			var sts Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var ro *jsonapi.ResourceObject
				ro, sts = readRequestObject(ctx)
				if is.Error(sts) {
					break
				}
				sts = ctx.Controller.AddItem(ctx, ro)
				if is.Error(sts) {
					break
				}
				rd.Data = ro.ResourceIdObject
			}
			return me.JsonMarshalHandler(ctx, sts)
		}),
	}
}

func (me *Api) WireUpdateItemRoute(e *echo.Echo, lc ListController, path string) *Route {
	return &Route{
		Route: e.PUT(path, func(ec echo.Context) error {
			var sts Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var id ItemId
				id, sts = apiworks.GetIdFromUrl(ctx, lc)
				if is.Error(sts) {
					break
				}
				var ro *jsonapi.ResourceObject
				ro, sts = readRequestObject(ctx)
				if is.Error(sts) {
					break
				}
				if ro.ResourceId != jsonapi.ResourceId(id) {
					sts = status.Fail(&status.Args{
						HttpStatus: http.StatusBadRequest,
						Message: fmt.Sprintf("id '%s' does not match ID segment in url '%s'",
							ro.ResourceId,
							ctx.Request().URL,
						),
					})
				}
				sts = ctx.Controller.UpdateItem(ctx, ro)
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
		}),
	}
}

func (me *Api) WireDeleteItemRoute(e *echo.Echo, lc ListController, path string) *Route {

	return &Route{
		Route: e.DELETE(path, func(ec echo.Context) error {
			var sts Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var id ItemId
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
		}),
	}
}

func (me *Api) GetItemUrl(ctx *Context, item ItemModeler) (u types.UrlTemplate, sts Status) {
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

func (me *Api) GetSelfUrl(ctx *Context) types.UrlTemplate {
	r := ctx.Request()
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, r.Host, me.GetSelfPath(ctx))
	return types.UrlTemplate(url)
}

func (me *Api) GetSelfPath(ctx *Context) types.UrlTemplate {
	return types.UrlTemplate(ctx.Request().RequestURI)
}

func (me *Api) GetContentType(ctx *Context) HttpHeaderValue {
	return jsonapi.ContentType + "; " + CharsetUTF8
}

func (me *Api) JsonMarshalHandler(ctx *Context, sts Status) Status {
	var _sts Status
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

	ctx.AddLink(SelfRelType, Link(me.GetSelfPath(ctx)))
	ctx.AddLink(SchemaDcRelType, DcSchema)
	ctx.AddLink(SchemaDcTermsRelType, DcTermsSchema)
	ctx.AddLink(SchemaGearboxApiRelationType, GearboxApiSchema)

	ctx.AddMeta(MetaDcCreator, GearboxApiIdentifier)
	ctx.AddMeta(MetaDcTermsIdentifier, me.GetSelfUrl(ctx))
	ctx.AddMeta(MetaDcLanguage, DefaultLanguage)
	ctx.AddMeta(MetaGearboxBaseurl, me.GetBaseUrl())

	if sts == nil {
		sts = status.Success("")
	}
	sts = ctx.SendResponse(sts)
	return sts
}

func (me *Api) AddController(controller ListController) (sts Status) {
	for range only.Once {
		getter, ok := controller.(BasepathGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "model does not implement BasepathGetter",
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
	return types.UrlTemplate(fmt.Sprintf(
		string(BaseUrlPattern),
		me.Port,
	))
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

func getResourceObject(rd *jsonapi.RootDocument) (ro *jsonapi.ResourceObject, sts Status) {
	ro, ok := rd.Data.(*jsonapi.ResourceObject)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: "root document does not contain a single resource object",
		})
	}
	ro.Renew()
	return ro, sts
}

func (me *Api) setItemData(ctx *Context, ro *jsonapi.ResourceObject, item ItemModeler, list List) (sts Status) {
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
		getter, ok := item.(ItemLinkMapGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "item does not implement ItemLinkMapGetter",
			})
			break
		}
		lm, sts := getter.GetItemLinkMap(ctx)
		if is.Error(sts) {
			break
		}
		lm.AddLink(SelfRelType, Link(su))
		ro.SetLinks(lm)
	}
	return sts
}

func (me *Api) SetRelationshipLinkMap(ctx *Context, item ItemModeler, ro *jsonapi.ResourceObject) (sts Status) {
	for range only.Once {
		lm, sts := ctx.RootDocumentor.GetDataRelationshipsLinkMap()
		if is.Error(sts) {
			break
		}
		name := item.GetType()

		for rel, link := range lm {

			rel = RelType(fmt.Sprintf("%s.%s", name, rel))

			ctx.RootDocumentor.AddLink(getQualifiedRelType(rel), link)
		}
	}
	return sts
}

func (me *Api) setListData(ctx *Context, data interface{}) (sts Status) {
	for range only.Once {
		getter, ok := data.(ListGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "data does not implement ListGetter",
			})
			break
		}
		var coll List
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

func (me *Api) getRouteName(ctx *Context) (name types.RouteName, sts Status) {
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

func getRootDocumentAndContext(ec echo.Context, lc ListController, rt types.ResponseType) (rd *jsonapi.RootDocument, ctx *Context) {
	rd = jsonapi.NewRootDocument(ec, rt)
	ctx = apiworks.NewContext(&ContextArgs{
		Contexter:      ec,
		RootDocumentor: rd,
		Controller:     lc,
	})
	return rd, ctx
}

func readRequestObject(ctx *Context) (ro *jsonapi.ResourceObject, sts Status) {
	defer closeRequestBody(ctx)
	for range only.Once {
		if ctx == nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("context not passed to '%s'", util.CurrentFunc()),
			})
			break
		}
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				HttpStatus: http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("unable to read body of '%s' request",
					ctx.Request().URL,
				),
			})
			break
		}
		ro = &jsonapi.ResourceObject{}
		err = ro.Unmarshal(body)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal body of '%s' request",
					ctx.Request().URL,
				),
				HttpStatus: http.StatusBadRequest,
			})
			break
		}
	}
	return ro, sts
}

func getQualifiedRelType(reltype RelType) RelType {
	rt := fmt.Sprintf(RelTypePattern, GearboxApiIdentifier, reltype)
	return RelType(rt)
}
