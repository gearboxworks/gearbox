package apimvc

import (
	"gearbox/apiworks"
	"gearbox/config"
	"gearbox/gearspec"
	"gearbox/jsonapi"
	"gearbox/types"
	"gearbox/util"
)

const BasedirModelType ItemType = "basedirs"

var NilBasedirModel = (*BasedirModel)(nil)
var _ ItemModeler = NilBasedirModel
var _ jsonapi.ResourceContainer = NilBasedirModel

type BasedirModelMap map[gearspec.Identifier]*BasedirModel
type BasedirModels []*BasedirModel

type BasedirModel struct {
	Nickname types.Nickname    `json:"nickname"`
	Basedir  types.AbsoluteDir `json:"basedir"`
	Model
}

func (me *BasedirModel) GetAttributeMap() (am apiworks.AttributeMap) {
	sm := util.StructMap(me, "json")
	am = make(apiworks.AttributeMap, len(sm))
	for k, v := range sm {
		am[Fieldname(k)] = v
	}
	return am
}

func (me *BasedirModel) ContainsResource() {}

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
