package apimodeler

type Context struct {
	Contexter
	RootDocumentor
	Controller ListController
}
type ContextArgs Context

func NewContext(args *ContextArgs) *Context {
	c := Context{}
	c = Context(*args)
	return &c
}
