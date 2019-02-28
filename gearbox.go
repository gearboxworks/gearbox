package gearbox

import (
	"context"
	"encoding/json"
	"fmt"
	"gearbox/dockerhub"
	"gearbox/host"
	"github.com/zserge/webview"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var Instance *Gearbox

type Gearbox struct {
	Config        *Config
	HostConnector host.Connector
	Stacks        StackMap
}

type GearboxArgs Gearbox

func (me *Gearbox) Initialize() {
	me.WriteAssetsToAdminWebRoot()
	me.Config.Initialize()
}

func NewGearbox(args *GearboxArgs) *Gearbox {
	if args.Config == nil {
		args.Config = NewConfig(args.HostConnector)
	}
	gb := Gearbox{
		HostConnector: args.HostConnector,
		Config:        args.Config,
		Stacks:        GetStackMap(),
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

func addCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		next.ServeHTTP(w, r)
	})
}
func shutdownServer(srv *http.Server) {
	fmt.Print("Shutting down Gearbox Admin Console web server...\n")
	err := srv.Shutdown(context.TODO())
	fmt.Print("Done.")
	if err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

//
// [1] https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
// [2] https://github.com/zserge/webview
//
func (me *Gearbox) Admin() {
	adminRootDir := me.HostConnector.GetAdminRootDir()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	srv := &http.Server{
		Handler: addCorsMiddleware(http.FileServer(http.Dir(adminRootDir))),
	}
	defer shutdownServer(srv)

	hostname := ln.Addr().String()

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}

	}()

	api := NewHostApi(me)
	apiJson := fmt.Sprintf(`{"host_api":"%s","vm_api":"%s"}`, api.Url(), api.Url())
	apiJsonFile := fmt.Sprintf("%s/api.json", adminRootDir)
	err = ioutil.WriteFile(apiJsonFile, []byte(apiJson), os.ModePerm)
	if err != nil {
		panic(err)
	}
	go api.Start()
	defer api.Stop()
	// See [2]
	wv := webview.New(webview.Settings{
		Title:     "Gearbox Admin Console",
		Height:    600,
		Width:     800,
		Resizable: true,
		URL:       fmt.Sprintf("http://%s/index.html", hostname),
		Debug:     true,
	})
	wv.Run()
}

func (me *Gearbox) WriteAssetsToAdminWebRoot() {
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

func (me *Gearbox) StartVm() {
	vm := &Vm{}
	err := vm.StartVm()
	if err != nil {
		panic(err)
	}
	return
}

func (me *Gearbox) RequestAvailableContainers(query ...*dockerhub.ContainerQuery) dockerhub.ContainerNames {
	var _query *dockerhub.ContainerQuery
	if len(query) == 0 {
		_query = &dockerhub.ContainerQuery{}
	} else {
		_query = query[0]
	}
	dh := dockerhub.DockerHub{}
	return dh.RequestAvailableContainerNames(_query)
}
