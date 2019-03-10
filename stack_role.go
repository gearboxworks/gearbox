package gearbox

const StackRoleHelpUrl = "https://docs.gearbox.works/roles/"

//
// Examples:
//
// 		wordpress/dbserver:1
//

type RoleName string

type StackRole struct {
	*Spec
}
type StackRoleArgs StackRole

type StackRoleDefaultGetter interface {
	DefaultHostGetter
	DefaultNamespaceGetter
}
type DefaultHostGetter interface {
	GetDefaultHost() string
}
type DefaultNamespaceGetter interface {
	GetDefaultNamespace() string
}
type DefaultRoleGetter interface {
	GetDefaultRole() string
}
var _ DefaultHostGetter = (*StackRole)(nil)
var _ DefaultNamespaceGetter = (*StackRole)(nil)

func (me *StackRole) GetDefaultHost() string {
	return "github.com"
}

func (me *StackRole) GetDefaultNamespace() string {
	return "gearboxworks"
}

