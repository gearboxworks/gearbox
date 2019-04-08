package ja

type Contexter interface {
	ParamGetter
	KeyValueGetter
	KeyValueSetter
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
