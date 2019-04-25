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
	HostDir  types.AbsoluteDir `json:"host_dir"`
	BoxDir   types.AbsoluteDir `json:"box_dir"`
}

func NewModelFromConfigBasedir(ctx *Context, bd *config.Basedir) (s *BasedirModel, sts Status) {
	s = &BasedirModel{
		Nickname: bd.Nickname,
		HostDir:  bd.HostDir,
		BoxDir:   bd.BoxDir,
	}
	return s, sts
}

func (me *BasedirModel) GetId() ItemId {
	return ItemId(me.Nickname)
}

func (me *BasedirModel) SetStackId(ItemId) Status {
	panic("implement me")
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
