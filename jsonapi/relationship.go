package ja

import (
	"gearbox/apimodeler"
)

type RelationshipType string

const (
	ToOneType  RelationshipType = "to-one"
	ToManyType RelationshipType = "to-many"
)

type RelationshipMap map[Fieldname]Relationship

type Relationship struct {
	Type    RelationshipType   `json:"-"`
	Meta    MetaMap            `json:"meta,omitempty"`
	Data    ResourceIdentifier `json:"data,omitempty"`
	LinkMap apimodeler.LinkMap `json:"links,omitempty"`
}
