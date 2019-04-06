package apimodeler

type ItemIds []ItemId
type ItemId string
type ItemType string

type FilterName string
type FilterLabel string
type FilterPath string
type FilterMap map[FilterPath]Filter
type Filters []Filter
type Critera interface{}
type Filter struct {
	Label            FilterLabel
	Path             FilterPath
	ItemFilter       func(Itemer) Itemer
	CollectionFilter func(Collection) Collection
}
