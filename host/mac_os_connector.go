package host

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"log"
)

const macOsUserDataPath = "Library/Application Support"
const macOsSuggestedProjectsPath = "Sites"

var MacOsConnectorType = (*MacOsConnector)(nil)

var _ Connector = MacOsConnectorType

type MacOsConnector struct{}

//
// See https://developer.apple.com/library/archive/documentation/FileManagement
//           /Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
//
//
func (me *MacOsConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s/%s/%s",
		me.GetUserHomeDir(),
		macOsUserDataPath,
		BundleIdentifier,
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
		macOsSuggestedProjectsPath,
	)
}

func (me *MacOsConnector) GetUserHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDir
}
