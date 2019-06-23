package global

import "gearbox/types"

const (
	DefaultAuthorityDomain            = "gearbox.works"
	DefaultOrgName                    = "gearboxworks"
	Brandname                         = "Gearbox"
	RelPrefix                         = "gearbox"
	UserDataPath           types.Path = ".gearbox"
)

const (
	ItemResponse types.ResponseType = "item"
	ListResponse types.ResponseType = "list"
)

type ShareableChoice string

const (
	NotShareable    ShareableChoice = "no"
	InStackSharable ShareableChoice = "instack"
	YesShareable    ShareableChoice = "yes"
)
