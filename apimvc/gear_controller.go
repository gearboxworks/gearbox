package apimvc

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"net/http"
	"reflect"
	"sort"
)

const GearControllerName types.RouteName = "gears"
const GearsBasepath types.Basepath = "/gears"

var NilGearController = (*GearController)(nil)
var _ ListController = NilGearController

type GearController struct {
	Controller
	Gearbox gearbox.Gearboxer
}

func NewGearController(gb gearbox.Gearboxer) *GearController {
	return &GearController{
		Gearbox: gb,
	}
}

func (me *GearController) GetRootObject() interface{} {
	return me.Gearbox
}

func (me *GearController) GetNilItem(ctx *Context) ItemModeler {
	return NilGearModel
}

func (me *GearController) GetRelatedFields() RelatedFields {
	return RelatedFields{}
}

func (me *GearController) GetName() types.RouteName {
	return GearControllerName
}

func (me *GearController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("http://example.com"),
	}, sts
}

func (me *GearController) GetBasepath() types.Basepath {
	return GearsBasepath
}

func (me *GearController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *GearController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		gbss, sts := me.Gearbox.GetServices()
		if is.Error(sts) {
			break
		}
		for _, gbs := range gbss {
			ns, sts := NewModelFromGear(ctx, gbs)
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

func (me *GearController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *GearController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
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

func (me *GearController) GetItem(ctx *Context, serviceid ItemId) (list ItemModeler, sts Status) {
	var ns *GearModel
	for range only.Once {
		gbns, sts := me.Gearbox.FindService(service.Identifier(serviceid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Gear '%s' not found", serviceid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewModelFromGear(ctx, gbns)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Gear '%s' found", serviceid)
	}
	return ns, sts
}

func (me *GearController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *GearController) GetFilterMap() FilterMap {
	return GetGearFilterMap()
}

func GetGearFilterMap() FilterMap {
	return FilterMap{}
}

func assertGear(item ItemModeler) (s *GearModel, sts Status) {
	s, ok := item.(*GearModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Gear: %v", item),
		})
	}
	return s, sts
}
