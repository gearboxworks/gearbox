package apiworks

var NilModel = (*Model)(nil)
var _ ItemModeler = NilModel

type Model struct {
}

func (me *Model) GetId() ItemId {
	panic("implement me")
}

func (me *Model) SetId(ItemId) Status {
	panic("implement me")
}

func (me *Model) GetType() ItemType {
	panic("implement me")
}

func (me *Model) GetAttributeMap() AttributeMap {
	panic("implement me")
}

func (me *Model) GetItem() (ItemModeler, Status) {
	return me, nil
}

func (me *Model) GetItemLinkMap(*Context) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("https://example.com"),
	}, sts
}

func (me *Model) GetRelatedItems(ctx *Context) (list List, sts Status) {
	return make(List, 0), sts
}
