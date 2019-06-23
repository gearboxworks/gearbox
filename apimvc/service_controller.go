package apimvc

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/http"
	"reflect"
	"sort"
)

const ServiceControllerName types.RouteName = "services"
const ServicesBasepath types.Basepath = "/services"

const OrgnameIdParam IdParam = "orgname"
const ProgramVersionIdParam IdParam = "progver"

var NilServiceController = (*ServiceController)(nil)
var _ ListController = NilServiceController

type ServiceController struct {
	Controller
	Gearbox gearbox.Gearboxer
}

func NewServiceController(gb gearbox.Gearboxer) *ServiceController {
	return &ServiceController{
		Gearbox: gb,
	}
}

func (me *ServiceController) GetNilItem(ctx *Context) ItemModeler {
	return NilServiceModel
}

func (me *ServiceController) GetRelatedFields() RelatedFields {
	return RelatedFields{}
}

func (me *ServiceController) GetName() types.RouteName {
	return ServiceControllerName
}

func (me *ServiceController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("http://example.com"),
	}, sts
}

func (me *ServiceController) GetBasepath() types.Basepath {
	return ServicesBasepath
}

func (me *ServiceController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *ServiceController) GetIdParams() IdParams {
	return IdParams{
		OrgnameIdParam,
		ProgramVersionIdParam,
	}
}

func (me *ServiceController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		gbss, sts := me.Gearbox.GetServices()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbss {
			ns, sts := NewModelFromGearsService(ctx, gbs)
			if is.Error(sts) {
				break
			}
			list = append(list, ns)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].GetId() < list[j].GetId()
		})
	}
	return list, sts
}

func (me *ServiceController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *ServiceController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
	for range only.Once {
		if len(filterPath) == 0 {
			filterPath = []FilterPath{NoFilterPath}
		}
		list, sts := me.GetList(ctx, filterPath[0])
		if is.Error(sts) {
			break
		}
		itemids = make(ItemIds, len(list))
		i := 0
		for _, item := range list {
			itemids[i] = ItemId(item.GetId())
			i++
		}
	}
	return itemids, sts
}

func (me *ServiceController) GetItem(ctx *Context, serviceid ItemId) (list ItemModeler, sts Status) {
	var ns *ServiceModel
	for range only.Once {
		gbns, sts := me.Gearbox.FindService(service.Identifier(serviceid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Service '%s' not found", serviceid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewModelFromGearsService(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Service '%s' found", serviceid)
	}
	return ns, sts
}

func (me *ServiceController) GetItemDetails(ctx *Context, itemid ItemId) (ItemModeler, Status) {
	return me.GetItem(ctx, itemid)
}

func (me *ServiceController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *ServiceController) GetFilterMap() FilterMap {
	return GetServiceFilterMap()
}

func GetServiceFilterMap() FilterMap {
	return FilterMap{}
}

func assertService(item ItemModeler) (s *ServiceModel, sts Status) {
	s, ok := item.(*ServiceModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Service: %v", item),
		})
	}
	return s, sts
}
