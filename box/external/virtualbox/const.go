package virtualbox

const (

	exitCodeOK        = "0"
	exitCodeMissingVm = "1"
	exitCodeCmdError  = "2"

	DefaultHostOnlyIp          = "192.168.42.1"
	DefaultHostOnlyNetmask     = "255.255.255.0"
	DefaultHostOnlyDhcpLowerIp = "192.168.42.10"
	DefaultHostOnlyDhcpUpperIp = "192.168.42.254"
)
