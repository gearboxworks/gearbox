package status

type HelpTypeMap map[HelpType]*string
type HelpType string

const (
	AllHelp HelpType = "all"
	ApiHelp HelpType = "api"
	CliHelp HelpType = "cli"
)
