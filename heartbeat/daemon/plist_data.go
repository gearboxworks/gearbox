package daemon

type PlistData struct {
	Label       string
	Program     string
	ProgramArgs []string
	Path        string
	KeepAlive   bool
	RunAtLoad   bool
	PidFile     string
}
