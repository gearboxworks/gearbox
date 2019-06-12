package daemon

import (
	"encoding/json"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"os"
	"path/filepath"
	"strings"
)


func ReadJsonConfig(f string) (*ServiceConfig, error) {

	var err error
	var c ServiceConfig
	var fh *os.File

	for range only.Once {
		if f == "" {
			err = messages.ProduceError(entity.DaemonEntityName, "Daemon service JSON file not defined")
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


func (me *Daemon) ParsePaths(sc ServiceConfig, i string) string {

	if me.OsPaths == nil {
		return i
	}

	strReplace := map[string]string {
		"{{.LocalDir}}":              me.OsPaths.LocalDir.String(),
		"{{.UserHomeDir}}":           me.OsPaths.UserHomeDir.String(),
		"{{.AdminRootDir}}":          me.OsPaths.AdminRootDir.String(),
		"{{.CacheDir}}":              me.OsPaths.CacheDir.String(),
		"{{.SuggestedBasedir}}":      me.OsPaths.SuggestedBasedir.String(),
		"{{.UserConfigDir}}":         me.OsPaths.UserConfigDir.String(),
		"{{.EventBrokerDir}}":        me.OsPaths.EventBrokerDir.String(),
		"{{.EventBrokerWorkingDir}}": me.OsPaths.EventBrokerWorkingDir.String(),
		"{{.EventBrokerLogDir}}":     me.OsPaths.EventBrokerLogDir.String(),
		"{{.EventBrokerEtcDir}}":     me.OsPaths.EventBrokerEtcDir.String(),
		"{{.Port}}":                  sc.autoPort,	// sc.UrlPtr.Port(),
		"{{.Host}}":                  sc.autoHost,	// sc.UrlPtr.Hostname(),
	}

	for k, v := range strReplace {
		i = strings.ReplaceAll(i, k, v)
	}

	return i
}


func (me *Daemon) ParseNetwork(sc ServiceConfig, i string) string {

	strReplace := map[string]string {
		"{{.Port}}":	sc.UrlPtr.Port(),
		"{{.Host}}":	sc.UrlPtr.Hostname(),
	}

	for k, v := range strReplace {
		i = strings.ReplaceAll(i, k, v)
	}

	return i
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
