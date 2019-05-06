package apimvc

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/http"
	"reflect"
	"sort"
)

const StackControllerName types.RouteName = "stacks"
const StacksBasepath types.Basepath = "/stacks"
const AuthorityIdParam IdParam = "authority"
const StacknameIdParam IdParam = "stackname"

const StackRolesField Fieldname = "stack_roles"

var NilStackController = (*StackController)(nil)
var _ ListController = NilStackController

type StackController struct {
	Controller
	Gearbox gearbox.Gearboxer
}

func NewStackController(gb gearbox.Gearboxer) *StackController {
	return &StackController{
		Gearbox: gb,
	}
}

func (me *StackController) GetNilItem(ctx *Context) ItemModeler {
	return NilNamedStackModel
}

func (me *StackController) GetRelatedFields() RelatedFields {
	return RelatedFields{
		&RelatedField{
			Fieldname:   StackRolesField,
			IncludeType: NamedStackType,
		},
	}
}

func (me *StackController) GetName() types.RouteName {
	return StackControllerName
}

func (me *StackController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("foobarbaz"),
	}, sts
}

func (me *StackController) GetBasepath() types.Basepath {
	return StacksBasepath
}

func (me *StackController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *StackController) GetIdParams() IdParams {
	return IdParams{
		AuthorityIdParam,
		StacknameIdParam,
	}
}

func (me *StackController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		gbnsm, sts := me.Gearbox.GetNamedStackMap()
		if is.Error(sts) {
			break
		}
		for _, gbns := range gbnsm {

			ns, sts := NewNamedStackModelFromGearsNamedStack(ctx, gbns)
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

func (me *StackController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *StackController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
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

func (me *StackController) AddItem(ctx *Context, item ItemModeler) (sts Status) {
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

func (me *StackController) UpdateItem(ctx *Context, item ItemModeler) (sts Status) {
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

func (me *StackController) DeleteItem(ctx *Context, stackid ItemId) (sts Status) {
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

func (me *StackController) GetItem(ctx *Context, stackid ItemId) (list ItemModeler, sts Status) {
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
		ns, sts = NewNamedStackModelFromGearsNamedStack(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
	}
	return ns, sts
}

func (me *StackController) GetItemDetails(ctx *Context, itemid ItemId) (ItemModeler, Status) {
	return me.GetItem(ctx, itemid)
}

func (me *StackController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *StackController) GetFilterMap() FilterMap {
	return GetStackFilterMap()
}

func (me *StackController) extractGearboxStack(ctx *Context, item ItemModeler) (gbs *gears.NamedStack, list List, sts Status) {
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

func GetStackFilterMap() FilterMap {
	return FilterMap{}
}

func assertStack(item ItemModeler) (s *NamedStackModel, sts Status) {
	s, ok := item.(*NamedStackModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}

func (me *StackController) getGearboxStackRoleMap() (gears.StackRoleMap, Status) {
	return me.Gearbox.GetStackRoleMap()
}
