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


func ReadJsonConfig(f string) (*ServiceConfig, error) {

	var err error
	var c ServiceConfig
	var fh *os.File

	for range only.Once {
		if f == "" {
			err = messages.ProduceError(DefaultEntityId, "Daemon service JSON file not defined")
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

	strReplace := map[string]string {
		"{{.GetLocalDir}}":			"/usr/local",
		"{{.GetUserHomeDir}}":		string(me.osSupport.GetUserHomeDir()),
		"{{.GetAdminRootDir}}":		string(me.osSupport.GetAdminRootDir()),
		"{{.GetCacheDir}}":			string(me.osSupport.GetCacheDir()),
		"{{.GetSuggestedBasedir}}":	string(me.osSupport.GetSuggestedBasedir()),
		"{{.GetUserConfigDir}}":	string(me.osSupport.GetUserConfigDir()),
		"{{.GetPort}}":				sc.Port.String(),
		"{{.GetHost}}":				sc.Host.String(),
	}

	for k, v := range strReplace {
		i = strings.ReplaceAll(i, k, v)
	}

	return i
}


func (me *Daemon) ParseNetwork(sc ServiceConfig, i string) string {

	strReplace := map[string]string {
		"{{.GetPort}}":	sc.Port.String(),
		"{{.GetHost}}":	sc.Host.String(),
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
