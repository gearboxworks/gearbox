// +build windows
package box

import (
	"os"
	"path/filepath"
)

// Hardcoded path to VBoxManage to fallback when it is not on path.
var VBOXMANAGE = filepath.Join(os.Getenv("PROGRAMFILES"), "Oracle", "VirtualBox", "VBoxManage.exe")
