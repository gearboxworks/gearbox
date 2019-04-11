package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
	"sort"
)

const GearspecsName types.RouteName = "gearspecs"
const GearspecsBasepath types.Basepath = "/gearspecs"
const RoleIdParam apimodeler.IdParam = "role"

var NilGearspecModel = (*GearspecModel)(nil)
var _ apimodeler.ApiModeler = NilGearspecModel

type GearspecModel struct {
	apimodeler.BaseModel
	Gearbox gearbox.Gearboxer
}

func NewGearspecModel(gb gearbox.Gearboxer) *GearspecModel {
	return &GearspecModel{
		Gearbox: gb,
	}
}
func (me *GearspecModel) CanAddItem(*apimodeler.Context) bool {
	return false
}

func (me *GearspecModel) GetName() types.RouteName {
	return GearspecsName
}

func (me *GearspecModel) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("foobarbaz"),
	}, sts
}

func (me *GearspecModel) GetBasepath() types.Basepath {
	return GearspecsBasepath
}

func (me *GearspecModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *GearspecModel) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
		RoleIdParam,
	}
}

func (me *GearspecModel) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	for range only.Once {
		gbgsrm, sts := me.Gearbox.GetGears().GetStackRoleMap()
		if is.Error(sts) {
			break
		}
		for _, gbgs := range gbgsrm {
			ns, sts := NewFromGearsGearspec(ctx, gbgs.Gearspec)
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

func (me *GearspecModel) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	return me.GetList(ctx, filterPath)
}

func (me *GearspecModel) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
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

func (me *GearspecModel) GetItem(ctx *apimodeler.Context, gearspecid apimodeler.ItemId) (list apimodeler.ApiItemer, sts status.Status) {
	var ns *Gearspec
	for range only.Once {
		gbgs, sts := me.Gearbox.GetGears().FindGearspec(gearspec.Identifier(gearspecid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Gearspec '%s' not found", gearspecid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewFromGearsGearspec(ctx, gbgs)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Gearspec '%s' found", gearspecid)
	}
	return ns, sts
}

func (me *GearspecModel) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.ApiItemer, status.Status) {
	return me.GetItem(ctx, itemid)
}

func (me *GearspecModel) FilterItem(in apimodeler.ApiItemer, filterPath apimodeler.FilterPath) (out apimodeler.ApiItemer, sts status.Status) {
	out = in
	return out, sts
}

func (me *GearspecModel) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertGearspec(item apimodeler.ApiItemer) (s *Gearspec, sts status.Status) {
	s, ok := item.(*Gearspec)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Gearspec: %v", item),
		})
	}
	return s, sts
}
