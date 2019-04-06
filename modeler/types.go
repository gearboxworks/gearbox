package modeler

type Item interface {
	ItemIdGetter
	ItemTypeGetter
	ItemGetter
}

type ItemIds []ItemId
type ItemId string
type ItemType string

type FilterName string
type FilterPath string
type FilterMap map[FilterPath]Filter
type Filters []Filter
type Critera interface{}
type Filter struct {
	Name   FilterName
	Path   FilterPath
	Filter func(Item) Item
}
