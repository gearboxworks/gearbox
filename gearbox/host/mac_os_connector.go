package host

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"log"
)

const bundleIdentifier = "works.gearbox.gearbox"
const userDataPath = "Library/Application Support"

var MacHostConnectorType = (*MacOsConnector)(nil)

var _ Connector = MacHostConnectorType

type MacOsConnector struct{}

//
// See https://developer.apple.com/library/archive/documentation/FileManagement
//           /Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
//
//
func (*MacOsConnector) GetUserDataDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/%s/%s", homeDir, userDataPath, bundleIdentifier)
}

func (me *MacOsConnector) GetWebRootDir() string {
	return fmt.Sprintf("%s/admin", me.GetUserDataDir())
}
