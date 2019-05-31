package messages

import (
	"github.com/gearboxworks/go-status"
)

func Debug(format string, a ...interface{}) {
	status.Success("DEBUG: " + format, a).Log()
}
