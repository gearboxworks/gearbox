package host

type Connector interface {
	GetWebRootDir() string
	GetUserDataDir() string
}
