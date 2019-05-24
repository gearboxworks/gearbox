package apimvc

import (
	"gearbox/apiworks"
	"gearbox/gearspec"
	"gearbox/jsonapi"
	"gearbox/types"
	"gearbox/util"
)

const DirectoryModelType ItemType = "directories"

var NilDirectoryModel = (*DirectoryModel)(nil)
var _ ItemModeler = NilDirectoryModel
var _ jsonapi.ResourceContainer = NilDirectoryModel

type DirectoryModelMap map[gearspec.Identifier]*DirectoryModel
type DirectoryModels []*DirectoryModel

type DirectoryModel struct {
	Directory types.AbsoluteDir `json:"directory"`
	Model
}

func (me *DirectoryModel) GetAttributeMap() (am apiworks.AttributeMap) {
	sm := util.StructMap(me, "json")
	am = make(apiworks.AttributeMap, len(sm))
	am["directory"] = sm["directory"]
	return am
}

func (me *DirectoryModel) ContainsResource() {}

func NewDirectoryModel(ctx *Context, dir types.AbsoluteDir) (s *DirectoryModel, sts Status) {
	s = &DirectoryModel{
		Directory: dir,
	}
	return s, sts
}

func (me *DirectoryModel) GetId() ItemId {
	return ItemId(me.Directory)
}

func (me *DirectoryModel) SetId(id ItemId) (sts Status) {
	me.Directory = types.AbsoluteDir(id)
	return sts
}

func (me *DirectoryModel) GetType() ItemType {
	return DirectoryModelType
}
