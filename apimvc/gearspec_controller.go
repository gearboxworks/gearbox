package apimvc

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

const GearspecControllerName types.RouteName = "gearspecs"
const GearspecsBasepath types.Basepath = "/gearspecs"
const RoleIdParam IdParam = "role"

var NilGearspecController = (*GearspecController)(nil)
var _ ListController = NilGearspecController

type GearspecController struct {
	Controller
	Gearbox gearbox.Gearboxer
}

func NewGearspecController(gb gearbox.Gearboxer) *GearspecController {
	return &GearspecController{
		Gearbox: gb,
	}
}
func (me *GearspecController) GetRelatedFields() RelatedFields {
	return RelatedFields{}
}

func (me *GearspecController) CanAddItem(*Context) bool {
	return false
}

func (me *GearspecController) GetName() types.RouteName {
	return GearspecControllerName
}

func (me *GearspecController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//StatusRelatedRelType: StatusLink("foobarbaz"),
	}, sts
}

func (me *GearspecController) GetBasepath() types.Basepath {
	return GearspecsBasepath
}

func (me *GearspecController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *GearspecController) GetIdParams() IdParams {
	return IdParams{
		AuthorityIdParam,
		StacknameIdParam,
		RoleIdParam,
	}
}

func (me *GearspecController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		gbgsrm, sts := me.Gearbox.GetGears().GetStackRoleMap()
		if is.Error(sts) {
			break
		}
		for _, gbgs := range gbgsrm {
			ns, sts := NewGearspecModelFromGearspecGearspec(ctx, gbgs.Gearspec)
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

func (me *GearspecController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *GearspecController) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts Status) {
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

func (me *GearspecController) GetItem(ctx *apimodeler.Context, gearspecid apimodeler.ItemId) (list apimodeler.ItemModeler, sts Status) {
	var ns *GearspecModel
	for range only.Once {
		gbgs, sts := me.Gearbox.GetGears().FindGearspec(gearspec.Identifier(gearspecid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Gearspec '%s' not found", gearspecid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewGearspecModelFromGearspecGearspec(ctx, gbgs)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Gearspec '%s' found", gearspecid)
	}
	return ns, sts
}

func (me *GearspecController) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.ItemModeler, Status) {
	return me.GetItem(ctx, itemid)
}

func (me *GearspecController) FilterItem(in apimodeler.ItemModeler, filterPath apimodeler.FilterPath) (out apimodeler.ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *GearspecController) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertGearspec(item ItemModeler) (s *GearspecModel, sts Status) {
	s, ok := item.(*GearspecModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Gearspec: %v", item),
		})
	}
	return s, sts
}
