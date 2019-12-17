package jsonapi

import (
	"gearbox/apiworks"
)

type (
	Fieldname       = apiworks.Fieldname
	Metaname        = apiworks.Metaname
	LinkImplementor = apiworks.LinkImplementor
	LinkMap         = apiworks.LinkMap
	Links           = apiworks.Links
	Link            = apiworks.Link
)

type (
	AttributeMap = apiworks.AttributeMap
	Attribute    = apiworks.Attribute
)

type Contexter interface {
	ParamGetter
	KeyValueGetter
	KeyValueSetter
}
type ParamGetter interface {
	Param(string) string
}
type KeyValueGetter interface {
	Get(string) interface{}
}
type KeyValueSetter interface {
	Set(string, interface{})
}

type LinkObject struct {
	Href    string  `json:"href"`
	MetaMap MetaMap `json:"meta"`
}

var _ apiworks.LinkImplementor = (*LinkObject)(nil)

func (*LinkObject) IdentifiesLink() {}

type MetaMap map[Metaname]interface{}

type ResponseError string
type Version string

type ResourceIds []ResourceId
type ResourceId = string
type ResourceTypes []ResourceType
type ResourceType = string

type ErrorId string
type ErrorCode string
type HttpStatus int
type JsonPointer string
type UrlParameter string
