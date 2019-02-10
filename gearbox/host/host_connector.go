package host

type Connector interface {
	GetAdminRootDir() string
	GetUserConfigDir() string
	GetSuggestedProjectRoot() string
	GetUserHomeDir() string
}
