package api

import "gearbox/status"

type UrlGetter interface {
	GetApiUrl(...RouteName) (UriTemplate, status.Status)
}
type UrlPathGetter interface {
	GetApiUrlPath(...RouteName) (UriTemplate, status.Status)
}

type ResponseDataGetter interface {
	GetResponseData() interface{}
}
