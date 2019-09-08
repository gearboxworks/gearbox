package vbmanage

const (
	VendorName     = "Oracle"
	ProductName    = "VirtualBox"
	ExecutableName = "VBoxManage"
)

const (
	StartVmCmd VbCmd = "startvm"
	ListCmd    VbCmd = "list"
)

const (
	exitCodeOK        = "0"
	exitCodeMissingVm = "1"
	exitCodeCmdError  = "2"
)
