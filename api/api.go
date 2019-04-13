package api

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/config"
	"gearbox/global"
	"gearbox/jsonapi"
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
					RootDocumenter: rd,
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
				RootDocumenter: rd,
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
				var item apimodeler.Itemer
				item, sts = ctlr.GetItemDetails(ctx, id)
				if is.Error(sts) {
					break
				}
				var list apimodeler.List
				list, sts = ctlr.GetRelatedItems(ctx, item)
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

func (me *Api) GetItemUrl(ctx *apimodeler.Context, item apimodeler.Itemer) (u types.UrlTemplate, sts status.Status) {
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

func (me *Api) GetContentType(ctx *apimodeler.Context) string {
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
	var err error
	var rootdoc *ja.RootDocument
	ec, ok := ctx.Contexter.(echo.Context)
	for range only.Once {
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "context does not implement echo.Context",
			})
			break
		}
		ec.Response().Header().Set(echo.HeaderContentType, me.GetContentType(ctx))

		rd := ctx.GetRootDocument()
		rootdoc, ok = rd.(*ja.RootDocument)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "context.RootDocument() does not implement ja.RootDocument",
			})
			break
		}

	}
	if rootdoc == nil {
		rootdoc = &ja.RootDocument{}
	}
	route := me.getRouteName(ec)
	if route == Rootname {
		rootdoc.AddLinks(me.GetRootLinkMap(ctx))
	}
	if rootdoc.GetResponseType() == global.ListResponse {
		rootdoc.AddLinks(me.GetListLinkMap(ctx))
	}
	if rootdoc.GetResponseType() == global.ItemResponse {
		rootdoc.AddLinks(me.GetItemLinkMap(ctx))
	}
	rootdoc.AddLinks(me.GetCommonLinkMap(ctx))

	rootdoc.MetaMap[apimodeler.MetaDcCreator] = GearboxApiIdentifier

	rootdoc.MetaMap[apimodeler.MetaDcTermsIdentifier] = me.GetSelfUrl(ctx)
	rootdoc.MetaMap[apimodeler.MetaDcLanguage] = apimodeler.DefaultLanguage

	rootdoc.LinkMap[apimodeler.SelfRelType] = apimodeler.Link(me.GetSelfPath(ctx))
	rootdoc.LinkMap[apimodeler.SchemaDcRelType] = apimodeler.DcSchema
	rootdoc.LinkMap[apimodeler.SchemaDcTermsRelType] = apimodeler.DcTermsSchema

	rootdoc.LinkMap[SchemaGearboxApiRelationType] = GearboxApiSchema
	rootdoc.MetaMap[MetaGearboxBaseurl] = me.GetBaseUrl()
	rootdoc.MetaMap[MetaGearboxApiSchema] = route

	if is.Error(sts) {
		ec.Response().Status = sts.HttpStatus()
		_ = rootdoc.SetError(sts)
	}
	err = ec.JSONPretty(ec.Response().Status, rootdoc, "   ")
	if err != nil {
		sts = status.Wrap(err, &status.Args{
			Message: fmt.Sprintf("error when sending output for '%s'",
				ec.Path(),
			),
		})
	}
	return sts
}

func (me *Api) AddController(controller apimodeler.ApiController) (sts status.Status) {
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

func (me *Api) setItemData(ctx *apimodeler.Context, ro *ja.ResourceObject, item apimodeler.Itemer, list apimodeler.List) (sts status.Status) {
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
		rootdoc, sts := assertRootDoc(ctx)
		if is.Error(sts) {
			break
		}

		sts = ro.SetRelatedItems(ctx, list)
		if is.Error(sts) {
			break
		}

		if rootdoc.ResponseType == global.ItemResponse {
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
		rootdoc, sts := assertRootDoc(ctx)
		if is.Error(sts) {
			break
		}
		for _, item := range coll {
			ro := ja.NewResourceObject()
			sts = me.setItemData(ctx, ro, item, nil)
			if is.Error(sts) {
				break
			}
			sts = rootdoc.AddResourceObject(ro)
			if is.Error(sts) {
				break
			}
		}
	}
	return sts
}

func (me *Api) getRouteName(ctx echo.Context) (name types.RouteName) {
	rts := me.Echo.Routes()
	path := ctx.Path()
	for _, rt := range rts {
		if rt.Path == path {
			name = types.RouteName(rt.Name)
			break
		}
	}
	return name
}

func assertRootDoc(ctx *apimodeler.Context) (rootdoc *ja.RootDocument, sts status.Status) {
	var ok bool
	rd := ctx.GetRootDocument()
	rootdoc, ok = rd.(*ja.RootDocument)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: "context.GetRootDocument() does not implement ja.RootDocument",
		})
	}
	return rootdoc, sts

}
func GetQualifiedRelType(reltype apimodeler.RelType) apimodeler.RelType {
	rt := fmt.Sprintf(apimodeler.RelTypePattern, GearboxApiIdentifier, reltype)
	return apimodeler.RelType(rt)
}
