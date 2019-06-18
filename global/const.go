package global

import "gearbox/types"

const (
	DefaultAuthority                    = "gearbox.works"
	DefaultOrgName                      = "gearboxworks"
	Brandname                           = "Gearbox"
	RelPrefix                           = "gearbox"
	UserDataPath     types.RelativePath = ".gearbox"
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
