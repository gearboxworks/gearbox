package api

import (
	"fmt"
	"gearbox/config"
	"gearbox/jsonapi" // @TODO Refactor this out to interface{}s
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"github.com/labstack/echo"
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

func (me *Api) SetParent(parent interface{}) {
	me.Parent = parent
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

func (me *Api) GetBaseUrl() (url types.UrlTemplate) {
	return types.UrlTemplate(fmt.Sprintf(
		string(BaseUrlPattern),
		me.Port,
	))
}

func (me *Api) SetRelationshipLinkMap(ctx *Context, item ItemModeler, ro *ResourceObject) (sts Status) {
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

func getQualifiedRelType(reltype RelType) RelType {
	rt := fmt.Sprintf(RelTypePattern, GearboxApiIdentifier, reltype)
	return RelType(rt)
}
