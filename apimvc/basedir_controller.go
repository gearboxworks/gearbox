package apimvc

import (
	"fmt"
	"gearbox/config"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"reflect"
	"sort"
)

const BasedirControllerName types.RouteName = "basedirs"
const BasedirsBasepath types.Basepath = "/basedirs"
const NicknameIdParam IdParam = "nickname"

var NilBasedirController = (*BasedirController)(nil)
var _ ListController = NilBasedirController

type BasedirController struct {
	Controller
	Config config.Configer
}

func NewBasedirController(cfg config.Configer) *BasedirController {
	return &BasedirController{
		Config: cfg,
	}
}

func (me *BasedirController) GetRelatedFields() RelatedFields {
	return RelatedFields{}
}

func (me *BasedirController) GetName() types.RouteName {
	return BasedirControllerName
}

func (me *BasedirController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("foobarbaz"),
	}, sts
}

func (me *BasedirController) GetBasepath() types.Basepath {
	return BasedirsBasepath
}

func (me *BasedirController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *BasedirController) GetIdParams() IdParams {
	return IdParams{
		NicknameIdParam,
	}
}

func (me *BasedirController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
	for range only.Once {
		bdm := me.Config.GetBasedirMap()
		for _, bd := range bdm {
			ns, sts := NewModelFromConfigBasedir(ctx, bd)
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

func (me *BasedirController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *BasedirController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
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

func (me *BasedirController) GetItem(ctx *Context, nickname ItemId) (list ItemModeler, sts Status) {
	var ns *BasedirModel
	for range only.Once {
		bd, sts := me.Config.FindBasedir(types.Nickname(nickname))
		if is.Error(sts) {
			break
		}
		ns, sts = NewModelFromConfigBasedir(ctx, bd)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Basedir '%s' found", nickname)
	}
	return ns, sts
}

func (me *BasedirController) GetItemDetails(ctx *Context, itemid ItemId) (ItemModeler, Status) {
	return me.GetItem(ctx, itemid)
}

func (me *BasedirController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *BasedirController) GetFilterMap() FilterMap {
	return GetBasedirFilterMap()
}

func (me *BasedirController) extractConfigBasedir(ctx *Context, item ItemModeler) (bd *config.Basedir, list List, sts Status) {
	var bdm *BasedirModel
	for range only.Once {
		list, sts = me.GetList(ctx)
		if is.Error(sts) {
			break
		}
		bdm, sts = assertBasedirModel(item)
		if is.Error(sts) {
			break
		}
		bd, sts = MakeConfigBasedir(me.Config, bdm)
	}
	return bd, list, sts
}

func GetBasedirFilterMap() FilterMap {
	return FilterMap{}
}

func MakeConfigBasedir(cfg config.Configer, bdm *BasedirModel) (bd *config.Basedir, sts Status) {
	bd = config.NewBasedir(bdm.HostDir, &config.BasedirArgs{
		Nickname: bdm.Nickname,
		BoxDir:   bdm.BoxDir,
	})
	return bd, sts
}

func assertBasedirModel(item ItemModeler) (bdm *BasedirModel, sts Status) {
	bdm, ok := item.(*BasedirModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a %T: %v",
				(*BasedirModel)(nil),
				item,
			),
		})
	}
	return bdm, sts
}

func (me *BasedirController) AddItem(ctx *Context, item ItemModeler) (sts Status) {
	//for range only.Once {
	//	var bd *config.Basedir
	//	bd, _, sts = me.extractConfigBasedir(ctx, item)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	sts = me.Config.AddBasedir(bd)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	sts = status.Success("Basedir '%s' added", bd.GetIdentifier())
	//	sts.SetHttpStatus(http.StatusCreated)
	//}
	return sts
}

func (me *BasedirController) UpdateItem(ctx *Context, item ItemModeler) (sts Status) {
	//for range only.Once {
	//	var gbs *config.Basedir
	//	gbs, _, sts = me.extractConfigBasedir(ctx, item)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	sts = me.Config.UpdateBasedir(gbs)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	sts = status.Success("Basedir '%s' updated", item.GetId())
	//	sts.SetHttpStatus(http.StatusNoContent)
	//}
	return sts

}

func (me *BasedirController) DeleteItem(ctx *Context, stackid ItemId) (sts Status) {
	//for range only.Once {
	//	sts := me.Config.DeleteBasedir(types.BasedirId(stackid))
	//	if status.IsError(sts) {
	//		break
	//	}
	//	sts = status.Success("Basedir '%s' found", stackid)
	//	sts.SetHttpStatus(http.StatusNoContent)
	//}
	return sts
}
