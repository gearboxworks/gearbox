package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
	"sort"
)

const StacksName types.RouteName = "stacks"
const StacksBasepath types.Basepath = "/stacks"
const AuthorityIdParam apimodeler.IdParam = "authority"
const StacknameIdParam apimodeler.IdParam = "stackname"

var NilStackController = (*StackController)(nil)
var _ apimodeler.ApiController = NilStackController

type StackController struct {
	apimodeler.Controller
	Gearbox gearbox.Gearboxer
}

func NewStackController(gb gearbox.Gearboxer) *StackController {
	return &StackController{
		Gearbox: gb,
	}
}

func (me *StackController) GetName() types.RouteName {
	return StacksName
}

func (me *StackController) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("foobarbaz"),
	}, sts
}

func (me *StackController) GetBasepath() types.Basepath {
	return StacksBasepath
}

func (me *StackController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *StackController) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
	}
}

func (me *StackController) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	for range only.Once {
		gbnsm, sts := me.Gearbox.GetNamedStackMap()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbnsm {
			ns, sts := NewModelFromGearsNamedStack(ctx, gbs)
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

func (me *StackController) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	return me.GetList(ctx, filterPath)
}

func (me *StackController) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
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

func (me *StackController) AddItem(ctx *apimodeler.Context, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.AddNamedStack(gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' added", gbs.GetIdentifier())
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func (me *StackController) UpdateItem(ctx *apimodeler.Context, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var gbs *gears.NamedStack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.UpdateNamedStack(gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' updated", item.GetId())
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts

}

func (me *StackController) DeleteItem(ctx *apimodeler.Context, stackid apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteNamedStack(types.StackId(stackid))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts
}

func (me *StackController) GetItem(ctx *apimodeler.Context, stackid apimodeler.ItemId) (list apimodeler.Itemer, sts status.Status) {
	var ns *NamedStackModel
	for range only.Once {
		gbns, sts := me.Gearbox.FindNamedStack(types.StackId(stackid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Stack '%s' not found", stackid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewModelFromGearsNamedStack(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
	}
	return ns, sts
}

func (me *StackController) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.Itemer, status.Status) {
	return me.GetItem(ctx, itemid)
}

func (me *StackController) FilterItem(in apimodeler.Itemer, filterPath apimodeler.FilterPath) (out apimodeler.Itemer, sts status.Status) {
	out = in
	return out, sts
}

func (me *StackController) GetFilterMap() apimodeler.FilterMap {
	return GetStackFilterMap()
}

func (me *StackController) extractGearboxStack(ctx *apimodeler.Context, item apimodeler.Itemer) (gbs *gears.NamedStack, list apimodeler.List, sts status.Status) {
	var ns *NamedStackModel
	for range only.Once {
		list, sts = me.GetList(ctx)
		if is.Error(sts) {
			break
		}
		ns, sts = assertStack(item)
		if is.Error(sts) {
			break
		}
		gbs, sts = MakeGearboxStack(me.Gearbox, ns)
	}
	return gbs, list, sts
}

func GetStackFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertStack(item apimodeler.Itemer) (s *NamedStackModel, sts status.Status) {
	s, ok := item.(*NamedStackModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}

func (me *StackController) getGearboxStackRoleMap() (gears.StackRoleMap, status.Status) {
	return me.Gearbox.GetStackRoleMap()
}
