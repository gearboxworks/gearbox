package jsonapi

type ErrorSource struct {
	JsonPointer  Link         `json:"pointer"`
	UrlParameter UrlParameter `json:"parameter"`
}
type ErrorSourceArgs ErrorSource

func NewErrorSource(args *ErrorSourceArgs) *ErrorSource {
	eo := ErrorSource{}
	eo = ErrorSource(*args)
	return &eo
}
