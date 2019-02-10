package host

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"log"
)

const bundleIdentifier = "works.gearbox.gearbox"
const userDataPath = "Library/Application Support"
const suggestedProjectsPath = "Sites"

var MacHostConnectorType = (*MacOsConnector)(nil)

var _ Connector = MacHostConnectorType

type MacOsConnector struct{}

//
// See https://developer.apple.com/library/archive/documentation/FileManagement
//           /Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
//
//
func (me *MacOsConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s/%s/%s",
		me.GetUserHomeDir(),
		userDataPath,
		bundleIdentifier,
	)
}

func (me *MacOsConnector) GetAdminRootDir() string {
	return fmt.Sprintf("%s/admin",
		me.GetUserConfigDir(),
	)
}

func (me *MacOsConnector) GetSuggestedProjectRoot() string {
	return fmt.Sprintf("%s/%s",
		me.GetUserHomeDir(),
		suggestedProjectsPath,
	)
}

func (me *MacOsConnector) GetUserHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDir
}
