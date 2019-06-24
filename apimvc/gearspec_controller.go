package apimvc

import (
	"fmt"
	"gearbox/apiworks"
	"gearbox/gearbox"
	"gearbox/gearspec"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"reflect"
	"sort"
)

const GearspecControllerName types.RouteName = "gearspecs"
const GearspecsBasepath types.Basepath = "/gearspecs"

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

func (me *GearspecController) GetNilItem(ctx *Context) ItemModeler {
	return NilGearspecModel
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

func (me *GearspecController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	fp := ""
	if len(filterPath) > 0 {
		fp = string(filterPath[0])
	}
	for range only.Once {
		gss, sts := me.Gearbox.GetGearRegistry().FilterGearspecByStack(types.StackId(fp))
		if is.Error(sts) {
			break
		}
		for _, sr := range gss {
			gs := gearspec.NewGearspec()
			sts := gs.Parse(sr.GearspecId)
			if is.Error(sts) {
				break
			}
			var gsm *GearspecModel
			gsm, sts = NewGearspecModelFromGearspecGearspec(ctx, gs)
			if is.Error(sts) {
				break
			}
			list = append(list, gsm)
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

func (me *GearspecController) GetListIds(ctx *apiworks.Context, filterPath ...apiworks.FilterPath) (itemids apiworks.ItemIds, sts Status) {
	for range only.Once {
		if len(filterPath) == 0 {
			filterPath = []apiworks.FilterPath{apiworks.NoFilterPath}
		}
		list, sts := me.GetList(ctx, filterPath[0])
		if is.Error(sts) {
			break
		}
		itemids = make(apiworks.ItemIds, len(list))
		i := 0
		for _, item := range list {
			itemids[i] = apiworks.ItemId(item.GetId())
			i++
		}
	}
	return itemids, sts
}

func (me *GearspecController) GetItem(ctx *apiworks.Context, gearspecid apiworks.ItemId) (item ItemModeler, sts Status) {
	var ns *GearspecModel
	for range only.Once {
		var gbgs *gearspec.Gearspec
		gbgs, sts = me.Gearbox.GetGearRegistry().FindGearspec(gearspec.Identifier(gearspecid))
		if is.Error(sts) {
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

func (me *GearspecController) FilterItem(in apiworks.ItemModeler, filterPath apiworks.FilterPath) (out apiworks.ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *GearspecController) GetFilterMap() apiworks.FilterMap {
	return apiworks.FilterMap{}
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
