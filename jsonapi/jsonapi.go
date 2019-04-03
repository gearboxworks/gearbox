package ja

import (
	"encoding/json"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
)

type RelationType string

const ContentType = "application/vnd.api+json"
const CharsetUTF8 = "charset=UTF-8"
const DcSchema Link = "http://purl.org/dc/elements/1.1/"
const DcTermsSchema Link = "http://purl.org/dc/terms/"
const DefaultLanguage = "en-US"

const (
	SchemaDcRelationType      RelationType = "schema.DC"
	SchemaDcTermsRelationType RelationType = "schema.DCTERMS"
	SelfRelationType          RelationType = "self"
	RelatedRelationType       RelationType = "related"
	FirstRelationType         RelationType = "first"
	LastRelationType          RelationType = "last"
	PrevRelationType          RelationType = "prev"
	NextRelationType          RelationType = "next"
	CollectionRelationType    RelationType = "collection"
	ItemRelationType          RelationType = "item"
)

type Metaname string

const (
	MetaDcLanguage        Metaname = "DC.language"
	MetaDcType            Metaname = "DC.type"
	MetaDcFormat          Metaname = "DC.format"
	MetaDcCreator         Metaname = "DC.creator"
	MetaDcExtent          Metaname = "DCTERMS.extent"
	MetaDcTermsIdentifier Metaname = "DCTERMS.identifier"
)

type MetaMap map[Metaname]interface{}

type Linker interface {
	IdentifiesLink()
}
type Link string
type LinkObject struct {
	Href    string  `json:"href"`
	MetaMap MetaMap `json:"meta"`
}

var _ Linker = (*LinkObject)(nil)
var _ Linker = Link("")

func (Link) IdentifiesLink()        {}
func (*LinkObject) IdentifiesLink() {}

type LinkMap map[RelationType]Linker
type Errors []error
type Version string

type JsonApi struct {
	Version Version `json:"version,omitempty"`
	MetaMap `json:"meta,omitempty"`
}

type DataSourcer interface {
	SourcesData()
}

type ResourceIdentifier interface {
	IdentifiesResource()
}
type ResourceContainer interface {
	ContainsResource()
}
type ResourceIds []ResourceId
type ResourceId string
type ResourceTypes []ResourceType
type ResourceType string

var _ ResourceIdentifier = (ResourceIdObjects)(nil)

func (ResourceIdObjects) IdentifiesResource() {}

type ResourceIdObjects []*ResourceIdObject

var _ ResourceIdentifier = (*ResourceIdObject)(nil)

func (*ResourceIdObject) IdentifiesResource() {}

type ResourceIdObject struct {
	ResourceId   `json:"id"`
	ResourceType `json:"type"`
	MetaMap      `json:"meta,omitempty"`
	LinkMap      `json:"links,omitempty"`
}

func NewResourceIdObject() *ResourceIdObject {
	rido := ResourceIdObject{
		MetaMap: make(MetaMap, 0),
		LinkMap: make(LinkMap, 0),
	}
	return &rido
}

type Fieldname string
type AttributeMap map[Fieldname]interface{}

var _ ResourceContainer = (ResourceObjects)(nil)

type ResourceObjects []*ResourceObject

func (ResourceObjects) ContainsResource() {}
func (me ResourceObjects) SetAttributes(attrs interface{}) (sts status.Status) {
	panic("Not yet implemented")
	return nil
}

func (me ResourceObjects) AppendResourceObject(ro *ResourceObject) (ResourceObjects, status.Status) {
	return append(me, ro), nil
}

type ResourceObjectAppender interface {
	AppendResourceObject(*ResourceObject) (ResourceObjects, status.Status)
}

type ResourceObjectsGetter interface {
	GetResourceObjects() ResourceObjects
}

type ResourceIdsSetter interface {
	SetIds(ResourceIds) status.Status
}

func (me ResourceObjects) SetIds(ids ResourceIds) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetId(ids[i])
		if is.Error(sts) {
			break
		}
	}
	return sts
}

type ResourceTypesSetter interface {
	SetTypes(ResourceTypes) status.Status
}

func (me ResourceObjects) SetTypes(types ResourceTypes) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetType(types[i])
		if is.Error(sts) {
			break
		}
	}
	return sts
}

var _ ResourceContainer = (*ResourceObject)(nil)

func (*ResourceObject) ContainsResource() {}

type ResourceObject struct {
	ResourceIdObject
	AttributeMap    `json:"attributes"`
	RelationshipMap `json:"relationships,omitempty"`
}

func NewResourceObject() *ResourceObject {
	ro := ResourceObject{
		AttributeMap:     make(AttributeMap, 0),
		RelationshipMap:  make(RelationshipMap, 0),
		ResourceIdObject: *NewResourceIdObject(),
	}
	return &ro
}
func (me *ResourceObject) SetId(id ResourceId) status.Status {
	me.ResourceId = id
	return nil
}

type ResourceIdGetter interface {
	GetId() ResourceId
}

func (me *ResourceObject) GetId() ResourceId {
	return me.ResourceId
}

type ResourceTypeGetter interface {
	GetType() ResourceType
}

func (me *ResourceObject) GetType() ResourceType {
	return me.ResourceType
}

type ResourceIdSetter interface {
	SetId(ResourceId) status.Status
}

type ResourceTypeSetter interface {
	SetType(ResourceType) status.Status
}

func (me *ResourceObject) SetType(_typ ResourceType) (sts status.Status) {
	me.ResourceType = _typ
	return nil
}

type AttributesSetter interface {
	SetAttributes(interface{}) status.Status
}

func (me *ResourceObject) SetAttributes(attrs interface{}) (sts status.Status) {
	var attrMap AttributeMap
	for range only.Once {
		b, err := json.Marshal(attrs)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: "unable to marshal attributes",
			})
			break
		}
		attrMap = make(AttributeMap, 0)
		err = json.Unmarshal(b, &attrMap)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: "unable to unmarshal attributes",
			})
			break
		}
		me.AttributeMap = attrMap
	}
	return sts
}
