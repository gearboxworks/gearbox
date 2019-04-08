package api

import "gearbox/types"

const (
	Port        = "9999"
	ServiceName = "Parent API"
	Version     = "0.1"
	DocsUrl     = "https://docs.gearbox.works/api"
)
const BaseUrlPattern = "http://127.0.0.1:%s"

const RequestContextKey = "api-request-context"

const (
	ItemResource    types.ResourceType = "item"
	ListResource    types.ResourceType = "list"
	MakeNewResource types.ResourceType = "new"
	UpdateResource  types.ResourceType = "update"
	DeleteResource  types.ResourceType = "delete"
)

const SelfResource types.RouteName = "self"
const LinksResource types.RouteName = "links"
