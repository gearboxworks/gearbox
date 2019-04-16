package apimodeler

import (
	"net/http"
)

type Contexter interface {
	ParamGetter
	KeyValueGetter
	KeyValueSetter
	RequestGetter
}

type RequestGetter interface {
	Request() *http.Request
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
