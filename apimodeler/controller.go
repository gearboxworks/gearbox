package apimodeler

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"net/http"
	"reflect"
)

const Basename = "controller"
const Basepath types.Basepath = "/"

var _ ApiController = (*Controller)(nil)

type ControllerMap map[types.Basepath]ApiController

type Controller struct {
	LinkMap  LinkMap
	Parent   *Controller
	Children ControllerMap
}

func (me *Controller) GetParent() ApiController {
	return me.Parent
}

func NewController() *Controller {
	return &Controller{
		LinkMap:  make(LinkMap, 0),
		Children: make(ControllerMap, 0),
	}
}

func (me *Controller) GetResourceUrlTemplate() (ut types.UrlTemplate) {
	for range only.Once {
		bp := me.GetBasepath()
		idt := GetIdTemplate(me)
		if idt == "" {
			ut = types.UrlTemplate(bp)
			break
		}
		ut = types.UrlTemplate(fmt.Sprintf("%s/%s", bp, idt))
	}
	return ut
}

func (me *Controller) AddLink(rel RelType, link LinkImplementor) {
	me.LinkMap[rel] = link
}
func (me *Controller) AddLinks(links LinkMap) {
	for rel, link := range links {
		me.AddLink(rel, link)
	}
}
func (me *Controller) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts status.Status) {
	return me.LinkMap, sts
}

func (me *Controller) CanAddItem(*Context) bool {
	return true
}

func (me *Controller) GetName() types.RouteName {
	return Basename
}

func (me *Controller) GetBasepath() types.Basepath {
	return Basepath
}

func (me *Controller) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *Controller) GetIdParams() IdParams {
	return IdParams{}
}

func (me *Controller) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts status.Status) {
	return list, sts
}

func (me *Controller) FilterList(ctx *Context, filterPath FilterPath) (list List, sts status.Status) {
	list = make(List, 0)
	return list, sts
}

func (me *Controller) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts status.Status) {
	itemids = make(ItemIds, 0)
	return itemids, sts
}

func (me *Controller) AddItem(ctx *Context, item Itemer) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *Controller) UpdateItem(ctx *Context, item Itemer) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *Controller) DeleteItem(ctx *Context, hostname ItemId) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *Controller) GetItem(ctx *Context, hostname ItemId) (item Itemer, sts status.Status) {
	return item, status.Success("Root found", hostname)
}

func (me *Controller) GetItemDetails(ctx *Context, hostname ItemId) (Itemer, status.Status) {
	return me.GetItem(ctx, hostname)
}

func (me *Controller) GetRelatedItems(ctx *Context, item Itemer) (list List, sts status.Status) {
	return make(List, 0), nil
}

func (me *Controller) FilterItem(in Itemer, filterPath FilterPath) (out Itemer, sts status.Status) {
	for range only.Once {
		if filterPath == NoFilterPath {
			out = in
			break
		}
		sts = status.Fail(&status.Args{
			Message:    fmt.Sprintf("filter '%s' not found", filterPath),
			HttpStatus: http.StatusBadRequest,
		})
		break
	}
	return out, sts
}

func (me *Controller) GetFilterMap() FilterMap {
	return FilterMap{}
}
