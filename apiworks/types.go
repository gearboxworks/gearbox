package apiworks

import "github.com/gearboxworks/go-status"

type (
	Status = status.Status
)

type (
	ItemIds  []ItemId
	ItemId   string
	ItemType string
)

type (
	FilterName  string
	FilterLabel string
	FilterPath  string
	FilterMap   map[FilterPath]Filter
	Filters     []Filter
)

type Filter struct {
	Label      FilterLabel
	Path       FilterPath
	ItemFilter func(ItemModeler) ItemModeler
	ListFilter func(List) List
}

type LinkImplementor interface {
	IdentifiesLink()
}
type Links []Link
type Link string

var _ LinkImplementor = Link("")
var _ LinkImplementor = Links{}

func (Links) IdentifiesLink() {}
func (Link) IdentifiesLink()  {}

type LinkMap map[RelType]LinkImplementor

func (me LinkMap) AddLink(reltype RelType, link LinkImplementor) {
	me[reltype] = link
}

type LinksSetter interface {
	SetLinks(LinkMap)
}

type Fieldname string
type RelatedFields []*RelatedField
type RelatedField struct {
	Fieldname   Fieldname
	LinkMap     LinkMap
	IncludeType ItemType
	Include     func(ItemModeler) bool
}

type HttpHeaderName string
type HttpHeaderValue string
type HttpResponseBody interface{}

type AttributeMap map[Fieldname]Attribute
type Attribute interface{}
