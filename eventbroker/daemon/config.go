package daemon

import (
	"encoding/json"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ReadJsonConfig(f string) (*ServiceConfig, error) {

	var err error
	var c ServiceConfig
	var fh *os.File

	for range only.Once {
		if f == "" {
			err = msgs.MakeError(entity.DaemonEntityName, "Daemon service JSON file not defined")
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

	if me.BaseDirs == nil {
		return i
	}

	strReplace := map[string]string{
		"{{.LocalDir}}":              me.BaseDirs.LocalDir,
		"{{.UserHomeDir}}":           me.BaseDirs.GetUserHomeDir(),
		"{{.AdminRootDir}}":          me.BaseDirs.GetAdminRootDir(),
		"{{.CacheDir}}":              me.BaseDirs.GetCacheDir(),
		"{{.ProjectBaseDir}}":        me.BaseDirs.GetProjectDir(),
		"{{.UserConfigDir}}":         me.BaseDirs.GetUserConfigDir(),
		"{{.EventBrokerDir}}":        me.BaseDirs.EventBrokerDir,
		"{{.EventBrokerWorkingDir}}": me.BaseDirs.EventBrokerWorkingDir,
		"{{.EventBrokerLogDir}}":     me.BaseDirs.EventBrokerLogDir,
		"{{.EventBrokerEtcDir}}":     me.BaseDirs.EventBrokerEtcDir,
		"{{.Port}}":                  sc.autoPort, // sc.UrlPtr.Port(),
		"{{.Host}}":                  sc.autoHost, // sc.UrlPtr.Hostname(),
		"{{.Platform}}":              runtime.GOOS + "_" + runtime.GOARCH,
	}

	for k, v := range strReplace {
		i = strings.ReplaceAll(i, k, v)
	}

	return i
}

func (me *Daemon) ParseNetwork(sc ServiceConfig, i string) string {

	strReplace := map[string]string{
		"{{.Port}}": sc.UrlPtr.Port(),
		"{{.Host}}": sc.UrlPtr.Hostname(),
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

func (c *ServiceConfig) SkipPlatform() (skip bool) {

	// Check platform.
	myPlatform := runtime.GOOS + "_" + runtime.GOARCH

	switch {
	case c.RunOnPlatform == "":
		skip = false

	case c.RunOnPlatform == myPlatform:
		skip = false

	default:
		skip = true
	}

	return
}
