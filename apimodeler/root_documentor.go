package apimodeler

import (
	"gearbox/status"
	"gearbox/types"
)

type RootDocumentor interface {
	ResponseTypeGetter
	RootDocumentGetter
	IncludedSetter
}
type ResponseTypeGetter interface {
	GetResponseType() types.ResponseType
}
type RootDocumentGetter interface {
	GetRootDocument() interface{}
}
type IncludedSetter interface {
	SetIncluded(*Context, List) status.Status
}
