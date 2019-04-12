package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
	"sort"
)

const ServicesName types.RouteName = "services"
const ServicesBasepath types.Basepath = "/services"

const OrgnameIdParam apimodeler.IdParam = "orgname"
const ServiceTypeIdParam apimodeler.IdParam = "svctype"
const ProgramVersionIdParam apimodeler.IdParam = "progver"

var NilServiceModel = (*ServiceModel)(nil)
var _ apimodeler.ApiModeler = NilServiceModel

type ServiceModel struct {
	apimodeler.BaseModel
	Gearbox gearbox.Gearboxer
}

func NewServiceModel(gb gearbox.Gearboxer) *ServiceModel {
	return &ServiceModel{
		Gearbox: gb,
	}
}

func (me *ServiceModel) GetName() types.RouteName {
	return ServicesName
}

func (me *ServiceModel) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("foobarbaz"),
	}, sts
}

func (me *ServiceModel) GetBasepath() types.Basepath {
	return ServicesBasepath
}

func (me *ServiceModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *ServiceModel) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		OrgnameIdParam,
		ServiceTypeIdParam,
		ProgramVersionIdParam,
	}
}

func (me *ServiceModel) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	for range only.Once {
		gbnsm, sts := me.Gearbox.GetServiceMap()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbnsm {
			ns, sts := NewFromGearsService(ctx, gbs)
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

func (me *ServiceModel) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	return me.GetList(ctx, filterPath)
}

func (me *ServiceModel) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
	for range only.Once {
		if len(filterPath) == 0 {
			filterPath = []apimodeler.FilterPath{apimodeler.NoFilterPath}
		}
		list, sts := me.GetList(ctx, filterPath[0])
		if is.Error(sts) {
			break
		}
		itemids = make(apimodeler.ItemIds, len(list))
		i := 0
		for _, item := range list {
			itemids[i] = apimodeler.ItemId(item.GetId())
			i++
		}
	}
	return itemids, sts
}

func (me *ServiceModel) GetItem(ctx *apimodeler.Context, serviceid apimodeler.ItemId) (list apimodeler.ApiItemer, sts status.Status) {
	var ns *Service
	for range only.Once {
		gbns, sts := me.Gearbox.FindService(types.ServiceId(serviceid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Service '%s' not found", serviceid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewFromGearsService(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Service '%s' found", serviceid)
	}
	return ns, sts
}

func (me *ServiceModel) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.ApiItemer, status.Status) {
	return me.GetItem(ctx, itemid)
}

func (me *ServiceModel) FilterItem(in apimodeler.ApiItemer, filterPath apimodeler.FilterPath) (out apimodeler.ApiItemer, sts status.Status) {
	out = in
	return out, sts
}

func (me *ServiceModel) GetFilterMap() apimodeler.FilterMap {
	return GetServiceFilterMap()
}

func GetServiceFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertService(item apimodeler.ApiItemer) (s *Service, sts status.Status) {
	s, ok := item.(*Service)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Service: %v", item),
		})
	}
	return s, sts
}
