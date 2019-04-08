package ja

import (
	"gearbox/apimodeler"
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

type MetaMap map[apimodeler.Metaname]interface{}
type Errors []error
type Version string

type ResourceIds []ResourceId
type ResourceId string
type ResourceTypes []ResourceType
type ResourceType string

type Fieldname string
type AttributeMap map[Fieldname]interface{}
