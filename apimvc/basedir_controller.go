package apimvc

import (
	"encoding/json"
	"fmt"
	"gearbox/config"
	"gearbox/jsonapi"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/http"
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

func (me *BasedirController) GetNilItem(ctx *Context) ItemModeler {
	return NilBasedirModel
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

func (me *BasedirController) getBasedirModelFromItem(item ItemModeler) (bdm *BasedirModel, sts Status) {
	for range only.Once {
		var ro *jsonapi.ResourceObject
		ro, sts = jsonapi.AssertResourceObject(item)
		if is.Error(sts) {
			break
		}
		if ro.GetType() != BasedirModelType {
			sts = status.Fail().
				SetHttpStatus(http.StatusBadRequest).
				SetMessage("invalid request type '%s'; should be '%s'",
					ro.GetType(),
					BasedirModelType,
				)
			break
		}
		var b []byte
		b, sts = ro.MarshalAttributeMap()
		if is.Error(sts) {
			break
		}
		err := json.Unmarshal(b, &bdm)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				HttpStatus: http.StatusBadRequest,
				Message: fmt.Sprintf("unable to marshal AttributeMap for Basedir '%s'",
					item.GetId(),
				),
			})
			break
		}
		sts = me.setBasedirModelId(bdm, item.GetId())
		if is.Error(sts) {
			break
		}
	}
	return bdm, sts
}

func (me *BasedirController) AddItem(ctx *Context, item ItemModeler) (sts Status) {
	for range only.Once {
		var bdm *BasedirModel
		bdm, sts = me.getBasedirModelFromItem(item)
		if is.Error(sts) {
			break
		}
		var bd *config.Basedir
		bd, sts = MakeConfigBasedir(bdm)
		if is.Error(sts) {
			break
		}
		bda := config.BasedirArgs(*bd)
		sts = me.Config.AddBasedir(&bda)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("base directory '%s' added", bd.Nickname).
			SetHttpStatus(http.StatusCreated).
			SetDetail("added the nicknamed '%s' directory '%s'",
				bda.Nickname,
				bda.Basedir,
			)
	}
	return sts
}

func (me *BasedirController) DeleteItem(ctx *Context, itemid ItemId) (sts Status) {
	return me.Config.DeleteBasedir(types.Nickname(itemid))
}

func (me *BasedirController) UpdateItem(ctx *Context, item ItemModeler) (sts Status) {
	for range only.Once {
		var bdm *BasedirModel
		bdm, sts = me.getBasedirModelFromItem(item)
		if is.Error(sts) {
			break
		}
		var bd *config.Basedir
		bd, sts = MakeConfigBasedir(bdm)
		if is.Error(sts) {
			break
		}
		sts = me.Config.UpdateBasedir(bd)
		if status.IsError(sts) {
			break
		}
	}
	return sts

}

func (me *BasedirController) setBasedirModelId(bdm *BasedirModel, itemid ItemId) (sts Status) {
	for range only.Once {
		nickname := bdm.GetId()
		if nickname == "" {
			sts = bdm.SetId(itemid)
			if is.Error(sts) {
				break
			}
		}
		if itemid == nickname {
			break
		}
		sts = status.Fail(&status.Args{
			HttpStatus: http.StatusUnprocessableEntity,
			Message: fmt.Sprintf("id '%s' does not match attributes.nickname '%s'",
				itemid,
				nickname,
			),
		})
	}
	return sts
}

func AssertBasedirModel(item ItemModeler) (bdm *BasedirModel, sts Status) {
	bdm, ok := item.(*BasedirModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a %T: %s",
				(*BasedirModel)(nil),
				item.GetId(),
			),
		})
	}
	return bdm, sts
}

func GetBasedirFilterMap() FilterMap {
	return FilterMap{}
}

func MakeConfigBasedir(bdm *BasedirModel) (bd *config.Basedir, sts Status) {
	bd = config.NewBasedir(bdm.Nickname, bdm.Basedir)
	return bd, sts
}
