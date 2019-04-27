package apimvc

import (
	"gearbox/apimodeler"
	"gearbox/jsonapi"
	"gearbox/status"
)

type (
	Status = status.Status
)

type (
	Context         = apimodeler.Context
	Controller      = apimodeler.Controller
	ControllerMap   = apimodeler.ControllerMap
	Fieldname       = apimodeler.Fieldname
	Filter          = apimodeler.Filter
	FilterMap       = apimodeler.FilterMap
	FilterPath      = apimodeler.FilterPath
	IdParam         = apimodeler.IdParam
	IdParams        = apimodeler.IdParams
	ItemId          = apimodeler.ItemId
	ItemIds         = apimodeler.ItemIds
	ItemModeler     = apimodeler.ItemModeler
	ItemType        = apimodeler.ItemType
	LinkImplementor = apimodeler.LinkImplementor
	LinkMap         = apimodeler.LinkMap
	List            = apimodeler.List
	ListController  = apimodeler.ListController
	RelatedField    = apimodeler.RelatedField
	RelatedFields   = apimodeler.RelatedFields
	RelType         = apimodeler.RelType
)

type (
	AttributeMap = jsonapi.AttributeMap
	Attribute    = jsonapi.Attribute
)
