package jsonapi

import (
	"gearbox/apiworks"
)

type RelationshipType string

const (
	ToOneType  RelationshipType = "to-one"
	ToManyType RelationshipType = "to-many"
)

type RelationshipMap map[Fieldname]*Relationship

type Relationship struct {
	Type    RelationshipType   `json:"-"`
	Meta    MetaMap            `json:"meta,omitempty"`
	Data    ResourceIdentifier `json:"data,omitempty"`
	LinkMap apiworks.LinkMap   `json:"links,omitempty"`
}

func NewRelationship() *Relationship {
	return &Relationship{
		Meta:    make(MetaMap, 0),
		LinkMap: make(apiworks.LinkMap, 0),
	}

}
