package host

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"log"
)

const winUserDataPath = ".gearbox"
const winSuggestedProjectsPath = "Gearbox Sites"

var WinConnectorType = (*WinConnector)(nil)

var _ Connector = WinConnectorType

type WinConnector struct{}

//
// See https://developer.apple.com/library/archive/documentation/FileManagement
//           /Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
//
//
func (me *WinConnector) GetUserConfigDir() string {
	return fmt.Sprintf("%s\\%s\\%s",
		me.GetUserHomeDir(),
		winUserDataPath,
		BundleIdentifier,
	)
}

func (me *WinConnector) GetAdminRootDir() string {
	return fmt.Sprintf("%s\\admin",
		me.GetUserConfigDir(),
	)
}

func (me *WinConnector) GetSuggestedProjectRoot() string {
	return fmt.Sprintf("%s\\%s",
		me.GetUserHomeDir(),
		winSuggestedProjectsPath,
	)
}

func (me *WinConnector) GetUserHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDir
}
