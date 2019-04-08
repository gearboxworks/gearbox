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

var _ Apier = (*Api)(nil)

type Api struct {
	Config    config.Configer
	Port      string
	Echo      *echo.Echo
	Parent    interface{}
	ModelsMap apimodeler.ModelsMap
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
		Config:    c.GetConfig(),
		Echo:      newConfiguredEcho(),
		Parent:    parent,
		ModelsMap: make(apimodeler.ModelsMap, 0),
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

func (me *Api) ConnectRoutes() {
	for _, models := range me.ModelsMap {

		// Copy to allow different values in closures
		ms := models

		e := me.Echo

		prefix := ms.GetRouteNamePrefix()

		var route *echo.Route

		path := string(ms.GetBasepath())

		// Collection Route
		route = e.GET(path, func(ctx echo.Context) error {
			rd := ja.NewRootDocument(ctx, ja.CollectionResponse)
			data, sts := ms.Self.GetCollection(ctx)
			if is.Success(sts) {
				sts = setCollectionData(ctx, rd, data)
			}
			return me.JsonMarshalHandler(rd, ctx, sts)
		})

		route.Name = fmt.Sprintf("%s-list", inflector.Pluralize(prefix))

		urlTemplate := string(ms.GetResourceUrlTemplate())

		// Single Item Route
		route = e.GET(urlTemplate, func(ctx echo.Context) error {
			var sts status.Status
			rd := ja.NewRootDocument(ctx, ja.DatasetResponse)
			for range only.Once {
				id, sts := ms.GetIdFromUrl(ctx)
				if is.Error(sts) {
					break
				}
				var ro *ja.ResourceObject
				ro, sts = getResourceObject(rd)
				if is.Error(sts) {
					break
				}
				var item apimodeler.Itemer
				item, sts = ms.Self.GetItem(ctx, id)
				if is.Error(sts) {
					break
				}
				sts = setItemData(ctx, ro, item)
				rd.Data = ro
			}
			return me.JsonMarshalHandler(rd, ctx, sts)
		})
		route.Name = fmt.Sprintf("%s-details", inflector.Singularize(prefix))
	}
}

func (me *Api) GetSelfUrl(ctx echo.Context) types.UrlTemplate {
	r := ctx.Request()
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, r.Host, me.GetSelfPath(ctx))
	return types.UrlTemplate(url)
}

func (me *Api) GetSelfPath(ctx echo.Context) types.UrlTemplate {
	return types.UrlTemplate(ctx.Request().RequestURI)
}

func (me *Api) GetContentType(ctx echo.Context) string {
	return ja.ContentType + "; " + ja.CharsetUTF8
}

const (
	GearboxApiSchema             ja.Link         = "https://docs.gearbox.works/api/schema/1.0/"
	SchemaGearboxApiRelationType ja.RelationType = "schema.GearboxAPI"
	MetaGearboxApiSchema         ja.Metaname     = "GearboxAPI.schema"
)

func (me *Api) JsonMarshalHandler(rootdoc *ja.RootDocument, ctx echo.Context, sts status.Status) status.Status {
	for range only.Once {
		var err error
		ctx.Response().Header().Set(echo.HeaderContentType, me.GetContentType(ctx))

		rootdoc.MetaMap[ja.MetaDcCreator] = global.Brandname

		rootdoc.MetaMap[ja.MetaDcTermsIdentifier] = me.GetSelfUrl(ctx)
		rootdoc.MetaMap[ja.MetaDcLanguage] = ja.DefaultLanguage

		rootdoc.LinkMap[ja.SelfRelationType] = ja.Link(me.GetSelfPath(ctx))
		rootdoc.LinkMap[ja.SchemaDcRelationType] = ja.DcSchema
		rootdoc.LinkMap[ja.SchemaDcTermsRelationType] = ja.DcTermsSchema

		rootdoc.LinkMap[SchemaGearboxApiRelationType] = GearboxApiSchema
		rootdoc.MetaMap[MetaGearboxApiSchema] = me.getRouteName(ctx)

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

func (me *Api) AddModels(models apimodeler.Modeler) (sts status.Status) {
	for range only.Once {
		getter, ok := models.(apimodeler.BasepathGetter)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "factory has no GetBasepath()",
			})
			break
		}
		path := getter.GetBasepath()
		me.ModelsMap[path] = apimodeler.NewModels(models)
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

func (me *Api) Url() string {
	return fmt.Sprintf(BaseUrlPattern, me.Port)
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

func setItemData(ctx apimodeler.Contexter, ro *ja.ResourceObject, item apimodeler.Itemer) (sts status.Status) {
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

func setCollectionData(ctx apimodeler.Contexter, rd *ja.RootDocument, data interface{}) (sts status.Status) {
	for range only.Once {
		getter, ok := data.(apimodeler.CollectionGetter)
		if !ok {
			break
		}
		var items apimodeler.Collection
		items, sts = getter.GetCollection(ctx)
		if is.Error(sts) {
			break
		}
		for _, item := range items {
			ro := ja.NewResourceObject()
			sts = setItemData(ctx, ro, item)
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
