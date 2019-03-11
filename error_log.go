package gearbox

import (
	"fmt"
	"gearbox/only"
	"os"
	"path/filepath"
)

type ErrorLog struct {
	Gearbox *Gearbox
}

func (me *ErrorLog) Write(b []byte) (nn int, err error) {
	for range only.Once {
		if me.Gearbox.IsDebug() {
			fmt.Print(string(b))
		}
		file := filepath.FromSlash(fmt.Sprintf("%s/error.log", me.Gearbox.HostConnector.GetUserConfigDir()))
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "Could not open '%s'", file)
			break
		}
		_, err = f.Write(b)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "Could not write '%s' to '%s'", string(b), file)
			break
		}
		nn = len(b)
		err = f.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "Could not close '%s'", file)
			break
		}
	}
	return nn, err
}
