package jsonapi

type ErrorObjects []*ErrorObject
type ErrorObject struct {
	ErrorId     ErrorId      `json:"id"`
	LinkMap     LinkMap      `json:"links"`
	HttpStatus  HttpStatus   `json:"status"`
	ErrorCode   ErrorCode    `json:"code"`
	Title       string       `json:"title"`
	Detail      string       `json:"detail"`
	ErrorSource *ErrorSource `json:"source"`
	MetaMap     MetaMap      `json:"meta,omitempty"`
}

type ErrorObjectArgs ErrorObject

func NewErrorObject(args *ErrorObjectArgs) *ErrorObject {
	eo := ErrorObject{}
	eo = ErrorObject(*args)
	return &eo
}
