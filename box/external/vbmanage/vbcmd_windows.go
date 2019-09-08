// +build windows

package vbmanage

import (
	"os"
	"path/filepath"
)

// Hardcoded path to VBoxManage to fallback when it is not on path.
var VBoxManagePath = filepath.Join(
	os.Getenv("PROGRAMFILES"),
	VendorName,
	ProductName,
	ExecutableName + ".exe",
)
