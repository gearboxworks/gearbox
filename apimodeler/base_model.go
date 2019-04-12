package apimodeler

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"net/http"
	"reflect"
)

const Basename = "base"
const Basepath types.Basepath = "/"

var _ ApiModeler = (*BaseModel)(nil)

type BaseModel struct {
	LinkMap LinkMap
}

func NewBaseModel() *BaseModel {
	return &BaseModel{
		LinkMap: make(LinkMap, 0),
	}
}
func (me *BaseModel) AddLink(rel RelType, link LinkImplementor) {
	me.LinkMap[rel] = link
}
func (me *BaseModel) AddLinks(links LinkMap) {
	for rel, link := range links {
		me.AddLink(rel, link)
	}
}
func (me *BaseModel) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts status.Status) {
	return me.LinkMap, sts
}

func (me *BaseModel) CanAddItem(*Context) bool {
	return true
}

func (me *BaseModel) GetName() types.RouteName {
	return Basename
}

func (me *BaseModel) GetBasepath() types.Basepath {
	return Basepath
}

func (me *BaseModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *BaseModel) GetIdParams() IdParams {
	return IdParams{}
}

func (me *BaseModel) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts status.Status) {
	return list, sts
}

func (me *BaseModel) FilterList(ctx *Context, filterPath FilterPath) (list List, sts status.Status) {
	list = make(List, 0)
	return list, sts
}

func (me *BaseModel) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts status.Status) {
	itemids = make(ItemIds, 0)
	return itemids, sts
}

func (me *BaseModel) AddItem(ctx *Context, item ApiItemer) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *BaseModel) UpdateItem(ctx *Context, item ApiItemer) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *BaseModel) DeleteItem(ctx *Context, hostname ItemId) (sts status.Status) {
	return status.Fail(&status.Args{
		Message:    "not supported",
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func (me *BaseModel) GetItem(ctx *Context, hostname ItemId) (item ApiItemer, sts status.Status) {
	return item, status.Success("Root found", hostname)
}

func (me *BaseModel) GetItemDetails(ctx *Context, hostname ItemId) (ApiItemer, status.Status) {
	return me.GetItem(ctx, hostname)
}

func (me *BaseModel) GetRelatedItems(ctx *Context, item ApiItemer) (list List, sts status.Status) {
	return make(List, 0), nil
}

func (me *BaseModel) FilterItem(in ApiItemer, filterPath FilterPath) (out ApiItemer, sts status.Status) {
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

func (me *BaseModel) GetFilterMap() FilterMap {
	return FilterMap{}
}
