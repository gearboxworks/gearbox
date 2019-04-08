package ja

import "gearbox/status"

type ResourceObjectAppender interface {
	AppendResourceObject(*ResourceObject) (ResourceObjects, status.Status)
}
type ResourceObjectsGetter interface {
	GetResourceObjects() ResourceObjects
}
type ResourceIdsSetter interface {
	SetIds(ResourceIds) status.Status
}
type ResourceTypesSetter interface {
	SetTypes(ResourceTypes) status.Status
}
type ResourceIdGetter interface {
	GetId() ResourceId
}
type ResourceTypeGetter interface {
	GetType() ResourceType
}
type ResourceIdSetter interface {
	SetId(ResourceId) status.Status
}
type ResourceTypeSetter interface {
	SetType(ResourceType) status.Status
}
type AttributesSetter interface {
	SetAttributes(interface{}) status.Status
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
