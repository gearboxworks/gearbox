package gears

import "gearbox/types"

const (
	//	RepoRawBaseUrl = "https://gears.gearbox.works"
	RepoRawBaseUrl = "https://raw.githubusercontent.com/gearboxworks/gearbox"
	JsonFilename   = "gears.json"
	JsonUrl        = RepoRawBaseUrl + "/master/assets/" + JsonFilename
	CacheKey       = "gears"
)

const (
	ServiceGear types.GearType = "service"
)

const ZeroString = ""
