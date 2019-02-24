package host

type Connector interface {
	GetAdminRootDir() string
	GetUserConfigDir() string
	GetSuggestedProjectRoot() string
	GetUserHomeDir() string
}

const BundleIdentifier = "works.gearbox.gearbox"
