package daemon

import (
	"encoding/json"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
	"os"
	"path/filepath"
	"strings"
)


func ReadJsonConfig(f string) (*CreateEntry, error) {

	var err error
	var c CreateEntry
	var fh *os.File

	for range only.Once {
		if f == "" {
			err = messages.ProduceError(defaultEntityId, "Daemon service JSON file not defined")
			break
		}

		fh, err = os.Open(f)
		if err != nil {
			break
		}
		defer fh.Close()

		r := json.NewDecoder(fh)
		err = r.Decode(&c)
		if err != nil {
			break
		}
	}

	if eblog.LogIfError(eblog.SkipNilCheck, err) {
		// Save last state.
		// me.State.Error = err
	}

	return &c, err
}


func (me *Daemon) ParsePaths(sc CreateEntry, i string) string {

	var o string

	i1 := strings.ReplaceAll(i, "{{.GetUserHomeDir}}", string(me.osSupport.GetUserHomeDir()))
	i2 := strings.ReplaceAll(i1, "{{.GetAdminRootDir}}", string(me.osSupport.GetAdminRootDir()))
	i3 := strings.ReplaceAll(i2, "{{.GetCacheDir}}", string(me.osSupport.GetCacheDir()))
	i4 := strings.ReplaceAll(i3, "{{.GetSuggestedBasedir}}", string(me.osSupport.GetSuggestedBasedir()))
	o = strings.ReplaceAll(i4, "{{.GetUserConfigDir}}", string(me.osSupport.GetUserConfigDir()))

	return o
}


func (me *Daemon) ParseNetwork(sc CreateEntry, i string) string {

	var o string

	i1 := strings.ReplaceAll(i, "{{.Port}}", sc.Port.String())
	o   = strings.ReplaceAll(i1, "{{.Host}}", sc.Host.String())

	return o
}


func (me *Daemon) CreateDirPaths(file string) error {

	var err error
	var dir string

	if file == "" {
		return nil
	}

	if !strings.HasPrefix(file, "/") {
		return nil
	}

	dir = filepath.Base(file)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	return err
}

// /Users/mick/.gearbox/admin/dist/eventbroker/logs/
