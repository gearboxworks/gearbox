package apimvc

import (
	"fmt"
	"gearbox/apiworks"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"net/http"
	"reflect"
	"sort"
)

const StackControllerName types.RouteName = "stacks"
const StacksBasepath types.Basepath = "/stacks"

const GearspecsField Fieldname = "gearspecs"

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
	return NilStackModel
}

func (me *StackController) GetRelatedFields() RelatedFields {
	return RelatedFields{
		&RelatedField{
			Fieldname:   GearspecsField,
			IncludeType: StackType,
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

func (me *StackController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		nss, sts := me.Gearbox.GetStacks()
		if is.Error(sts) {
			break
		}
		for _, ns := range nss {
			ns, sts := NewStackModelFromGearsStack(ctx, ns)
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

func (me *StackController) AddItem(ctx *Context, item ItemModeler) (im ItemModeler, sts Status) {
	for range only.Once {
		var gbs *gears.Stack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.AddStack(gbs)
		if status.IsError(sts) {
			break
		}
		im, sts = NewStackModelFromGearsStack(ctx, gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' added", gbs.GetIdentifier())
		_ = sts.SetHttpStatus(http.StatusCreated)
	}
	return im, sts
}

func (me *StackController) UpdateItem(ctx *apiworks.Context, item apiworks.ItemModeler) (modeler apiworks.ItemModeler, sts status.Status) {
	for range only.Once {
		var gbs *gears.Stack
		gbs, _, sts = me.extractGearboxStack(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.UpdateStack(gbs)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' updated", item.GetId())
		_ = sts.SetHttpStatus(http.StatusNoContent)
	}
	return nil, sts

}

func (me *StackController) DeleteItem(ctx *Context, stackid ItemId) (sts Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteStack(types.StackId(stackid))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
		_ = sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts
}

func (me *StackController) GetItem(ctx *Context, stackid ItemId) (list ItemModeler, sts Status) {
	var ns *StackModel
	for range only.Once {
		gbns, sts := me.Gearbox.FindStack(types.StackId(stackid))
		if is.Error(sts) {
			break
		}
		ns, sts = NewStackModelFromGearsStack(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Stack '%s' found", stackid)
	}
	return ns, sts
}

func (me *StackController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *StackController) GetFilterMap() FilterMap {
	return GetStackFilterMap()
}

func (me *StackController) extractGearboxStack(ctx *Context, item ItemModeler) (gbs *gears.Stack, list List, sts Status) {
	var ns *StackModel
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

func assertStack(item ItemModeler) (s *StackModel, sts Status) {
	s, ok := item.(*StackModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Stack: %v", item),
		})
	}
	return s, sts
}
