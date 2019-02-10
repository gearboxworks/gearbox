package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/gearbox/host"
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	AdminUpdater  func()
}

func (me *Gearbox) Initialize() {
	me.WriteAdminAssetsToWebRoot()
	me.Config.Initialize()
}

func NewGearbox(hc host.Connector) *Gearbox {
	gb := Gearbox{
		HostConnector: hc,
		Config:        NewConfig(hc),
	}
	return &gb
}

func (me *Gearbox) GetProjects() string {
	j, err := json.Marshal(me.Config.Projects)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

//
//
// [1] https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
// [2] https://github.com/zserge/webview
//
func (me *Gearbox) Admin() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		// See [1]
		http.Handle("/", http.FileServer(http.Dir(me.HostConnector.GetAdminRootDir())))
		log.Fatal(http.Serve(ln, nil))
	}()

	// See [2]
	wv := webview.New(webview.Settings{
		Title:     "Gearbox Admin Console",
		Height:    600,
		Width:     800,
		Resizable: true,
		URL:       fmt.Sprintf("http://%s/index.html", ln.Addr().String()),
		Debug:     true,
	})
	wv.Dispatch(func() {
		me.AdminUpdater, err = wv.Bind("gearbox", &Bridge{
			Webview: wv,
			Gearbox: me,
		})
		if err != nil {
			log.Fatal(err)
		}
		//err = wv.Eval(string(MustAsset("admin/js/vue.js")))
		//if err != nil {
		//	log.Fatal(err)
		//}
	})
	wv.Run()
}

func (me *Gearbox) WriteAdminAssetsToWebRoot() {
	hc := me.HostConnector
	if hc == nil {
		log.Fatal("Gearbox has no host connector. (End users should never see this; it is a programming error.)")
	}
	for _, afn := range AssetNames() {
		err := RestoreAsset(hc.GetUserConfigDir(), afn)
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not restore asset '%s/%s'",
				hc.GetUserConfigDir(),
				afn,
			))
		}
	}
}

func (me *Gearbox) AddProjectRoot(dir string) {
	pr := NewProjectRoot(me.Config.VmProjectRoot, dir)
	me.Config.ProjectRoots = append(me.Config.ProjectRoots, pr)
	me.Config.LoadProjectsAndWrite()
}
