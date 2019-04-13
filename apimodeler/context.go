package apimodeler

type Context struct {
	Contexter
	RootDocumenter
	Controller ApiController
}
type ContextArgs Context

func NewContext(args *ContextArgs) *Context {
	c := Context{}
	c = Context(*args)
	return &c
}
