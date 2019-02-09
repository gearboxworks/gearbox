package gearbox

import (
	"fmt"
	"gearbox/gearbox/host"
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
)

type Admin struct {
	Host host.Connector
}

func (me *Admin) Run() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		// See https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
		http.Handle("/", http.FileServer(http.Dir(me.Host.GetWebRootDir())))
		log.Fatal(http.Serve(ln, nil))
	}()

	//
	// See https://github.com/zserge/webview
	//
	w := webview.New(webview.Settings{
		Title:     "Gearbox Admin Console",
		Height:    600,
		Width:     800,
		Resizable: true,
		URL:       fmt.Sprintf("http://%s/index.html", ln.Addr().String()),
		Debug:     true,
	})
	////updaterFunc, err := w.Bind("admin", me)
	//_, err := w.Bind("admin", me)
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = w.Eval(string(MustAsset("admin/js/vue.js")))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	w.Run()
}

func WriteAdminAssetsToWebRoot(c host.Connector) {
	for _, afn := range AssetNames() {
		err := RestoreAsset(c.GetUserDataDir(), afn)
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not restore asset '%s/%s'", c.GetUserDataDir(), afn))
		}
	}

}
