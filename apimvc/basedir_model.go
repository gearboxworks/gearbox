package apimvc

import (
	"gearbox/config"
	"gearbox/gearspec"
	"gearbox/types"
)

const BasedirModelType ItemType = "basedir"

var NilBasedirModel = (*BasedirModel)(nil)
var _ ItemModeler = NilBasedirModel

type BasedirModelMap map[gearspec.Identifier]*BasedirModel
type BasedirModels []*BasedirModel

type BasedirModel struct {
	Nickname types.Nickname    `json:"nickname"`
	Basedir  types.AbsoluteDir `json:"basedir"`
}

func NewModelFromConfigBasedir(ctx *Context, bd *config.Basedir) (s *BasedirModel, sts Status) {
	s = &BasedirModel{
		Nickname: bd.Nickname,
		Basedir:  bd.Basedir,
	}
	return s, sts
}

func (me *BasedirModel) GetId() ItemId {
	return ItemId(me.Nickname)
}

func (me *BasedirModel) SetId(id ItemId) (sts Status) {
	me.Nickname = types.Nickname(id)
	return sts
}

func (me *BasedirModel) GetType() ItemType {
	return BasedirModelType
}

func (me *BasedirModel) GetItem() (ItemModeler, Status) {
	return me, nil
}

func (me *BasedirModel) GetItemLinkMap(*Context) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("https://example.com"),
	}, sts
}

func (me *BasedirModel) GetRelatedItems(ctx *Context) (list List, sts Status) {
	return make(List, 0), sts
}
