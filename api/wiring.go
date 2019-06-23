package api

import (
	"encoding/json"
	"fmt"
	"gearbox/apiworks"
	"gearbox/global"
	"gearbox/jsonapi"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

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

		me.WireHeadRoute(
			me.Echo,
			ctlr,
			itempath,
		).SetSingularName("head-%s", prefix)

		me.WireGetItemRoute(
			me.Echo,
			ctlr,
			itempath,
		).SetSingularName("get-%s", prefix)

		me.WireAddItemRoute(
			me.Echo,
			ctlr,
			prefix,
			listpath,
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
				err = me.Handler(ctx, sts)
			}
			return err
		}),
	}
}

func (me *Api) WireHeadRoute(e *echo.Echo, lc ListController, path string) *Route {

	return &Route{
		Route: e.HEAD(path, func(ec echo.Context) error {
			var sts status.Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var id ItemId
				id, sts = apiworks.GetIdFromUrl(ctx, lc)
				if is.Error(sts) {
					break
				}
				var ro *ResourceObject
				ro, sts = getResourceObject(rd)
				if is.Error(sts) {
					break
				}
				var item ItemModeler
				item, sts = lc.GetItem(ctx, id) //@TODO Maybe make this lighter weight for HEAD request
				if is.Error(sts) {
					break
				}
				sts = me.setItemData(ctx, ro, item, nil)
				rd.Data = ro
			}
			return me.Handler(ctx, sts)
		}),
	}
}

func (me *Api) WireGetItemRoute(e *echo.Echo, lc ListController, path string) *Route {

	return &Route{
		Route: e.GET(path, func(ec echo.Context) error {
			var sts status.Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var id ItemId
				id, sts = apiworks.GetIdFromUrl(ctx, lc)
				if is.Error(sts) {
					break
				}
				var ro *ResourceObject
				ro, sts = getResourceObject(rd)
				if is.Error(sts) {
					break
				}
				var item ItemModeler
				item, sts = lc.GetItem(ctx, id)
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
			return me.Handler(ctx, sts)
		}),
	}
}

func (me *Api) WireAddItemRoute(e *echo.Echo, lc ListController, prefix, path string) *Route {
	return &Route{
		Route: e.POST(path, func(ec echo.Context) error {
			var sts Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var ro *ResourceObject
				ro, sts = readRequestObject(ctx)
				if is.Error(sts) {
					break
				}
				var item ItemModeler
				item, sts = ctx.Controller.AddItem(ctx, ro)
				if is.Error(sts) {
					break
				}
				rd.Data = NewResourceObjectFromItem(item)
			}
			return me.Handler(ctx, sts)
		}),
	}
}

func (me *Api) WireUpdateItemRoute(e *echo.Echo, lc ListController, path string) *Route {
	return &Route{
		Route: e.PATCH(path, func(ec echo.Context) error {
			var sts Status
			rd, ctx := getRootDocumentAndContext(ec, lc, global.ItemResponse)
			for range only.Once {
				var id ItemId
				id, sts = apiworks.GetIdFromUrl(ctx, lc)
				if is.Error(sts) {
					break
				}
				var ro *ResourceObject
				ro, sts = readRequestObject(ctx)
				if is.Error(sts) {
					break
				}
				if ro.GetId() != id {
					sts = status.Fail(&status.Args{
						HttpStatus: http.StatusBadRequest,
						Message: fmt.Sprintf("id '%s' does not match ID segment in url '%s'",
							ro.GetId(),
							ctx.Request().URL,
						),
					})
					break
				}
				var item ItemModeler
				item, sts = ctx.Controller.UpdateItem(ctx, ro)
				if is.Error(sts) {
					break
				}
				rd.Data = NewResourceObjectFromItem(item)
			}
			return me.Handler(ctx, sts)
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
			return me.Handler(ctx, sts)
		}),
	}
}

func (me *Api) setItemData(ctx *Context, ro *ResourceObject, item ItemModeler, list List) (sts Status) {
	for range only.Once {

		sts = ro.SetId(item.GetId())
		if is.Error(sts) {
			break
		}

		sts = ro.SetType(item.GetType())
		if is.Error(sts) {
			break
		}

		if ctx.Contexter.Request().Method == http.MethodHead {
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

func NewResourceObjectFromItem(item ItemModeler) *ResourceObject {
	ro := ResourceObject{}
	ro.ResourceId = jsonapi.ResourceId(item.GetId())
	ro.ResourceType = jsonapi.ResourceType(item.GetType())
	ro.AttributeMap = item.GetAttributeMap()
	return &ro
}

func getResourceObject(rd *jsonapi.RootDocument) (ro *ResourceObject, sts Status) {
	ro, ok := rd.Data.(*ResourceObject)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: "root document does not contain a single resource object",
		})
	}
	ro.Renew()
	return ro, sts
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

func readRequestObject(ctx *Context) (ro *ResourceObject, sts Status) {
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
		do := &jsonapi.DataObject{}
		err = json.Unmarshal(body, &do)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal body of '%s' request",
					ctx.Request().URL,
				),
				HttpStatus: http.StatusBadRequest,
			})
			break
		}
		ro = do.Data
	}
	return ro, sts
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
