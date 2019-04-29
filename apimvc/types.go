package apimvc

import (
	"gearbox/apiworks"
	"gearbox/jsonapi"
	"gearbox/status"
)

type (
	Status = status.Status
)

type (
	Context         = apiworks.Context
	Controller      = apiworks.Controller
	ControllerMap   = apiworks.ControllerMap
	Fieldname       = apiworks.Fieldname
	Filter          = apiworks.Filter
	FilterMap       = apiworks.FilterMap
	FilterPath      = apiworks.FilterPath
	IdParam         = apiworks.IdParam
	IdParams        = apiworks.IdParams
	ItemId          = apiworks.ItemId
	ItemIds         = apiworks.ItemIds
	ItemModeler     = apiworks.ItemModeler
	ItemType        = apiworks.ItemType
	LinkImplementor = apiworks.LinkImplementor
	LinkMap         = apiworks.LinkMap
	List            = apiworks.List
	ListController  = apiworks.ListController
	RelatedField    = apiworks.RelatedField
	RelatedFields   = apiworks.RelatedFields
	RelType         = apiworks.RelType
)

type (
	AttributeMap = jsonapi.AttributeMap
	Attribute    = jsonapi.Attribute
)
