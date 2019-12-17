package global

import "gearbox/types"

const (
	DefaultAuthority            = "gearbox.works"
	DefaultOrgName              = "gearboxworks"
	Brandname                   = "Gearbox"
	RelPrefix                   = "gearbox"
	UserDataPath     types.Path = ".gearbox"
)

const (
	WindowsSuggestedBasedir types.Dir  = "Gearbox Sites"
	WindowsAdminPath        types.Path = "admin\\dist"
	LinuxSuggestedBasedir   types.Dir  = "projects"
	MacOsSuggestedBasedir   types.Dir  = "Sites"
	NixAdminPath            types.Path = "admin/dist"
)

const (
	ItemResponse types.ResponseType = "item"
	ListResponse types.ResponseType = "list"
)

type ShareableChoices string

const (
	NotShareable    ShareableChoices = "no"
	InStackSharable ShareableChoices = "instack"
	YesShareable    ShareableChoices = "yes"
)
