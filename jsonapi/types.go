package jsonapi

import (
	"gearbox/apimodeler"
)

type (
	Fieldname       = apimodeler.Fieldname
	Metaname        = apimodeler.Metaname
	LinkImplementor = apimodeler.LinkImplementor
	LinkMap         = apimodeler.LinkMap
	Links           = apimodeler.Links
	Link            = apimodeler.Link
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

var _ apimodeler.LinkImplementor = (*LinkObject)(nil)

func (*LinkObject) IdentifiesLink() {}

type MetaMap map[Metaname]interface{}

type ResponseError string
type Version string

type ResourceIds []ResourceId
type ResourceId string
type ResourceTypes []ResourceType
type ResourceType string

type AttributeMap map[Fieldname]interface{}
type Attribute interface{}

type ErrorId string
type ErrorCode string
type HttpStatus int
type JsonPointer string
type UrlParameter string
