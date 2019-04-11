package apimodeler

type ItemIds []ItemId
type ItemId string
type ItemType string

type FilterName string
type FilterLabel string
type FilterPath string
type FilterMap map[FilterPath]Filter
type Filters []Filter
type Filter struct {
	Label      FilterLabel
	Path       FilterPath
	ItemFilter func(ApiItemer) ApiItemer
	ListFilter func(List) List
}

type LinkImplementor interface {
	IdentifiesLink()
}
type Link string

var _ LinkImplementor = Link("")

func (Link) IdentifiesLink() {}

type LinkMap map[RelType]LinkImplementor

func (me LinkMap) AddLink(reltype RelType, link LinkImplementor) {
	me[reltype] = link
}

type LinksSetter interface {
	SetLinks(LinkMap)
}
