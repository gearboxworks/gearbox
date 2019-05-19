package api

import (
	"gearbox/global"
	"gearbox/jsonapi"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/labstack/echo"
)

func (me *Api) Handler(ctx *Context, sts Status) Status {
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
