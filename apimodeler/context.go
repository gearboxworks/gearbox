package apimodeler

type Context struct {
	Contexter
	RootDocumentGetter
	Models *Models
}
type ContextArgs Context

func NewContext(args *ContextArgs) *Context {
	c := Context{}
	c = Context(*args)
	return &c
}
