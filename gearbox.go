package gearbox

import (
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

//
//
// [1] https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
// [2] https://github.com/zserge/webview
//
func (me *Gearbox) Admin() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	adminRootDir := me.HostConnector.GetAdminRootDir()

	go func() {
		// See [1]
		http.Handle("/", http.FileServer(http.Dir(adminRootDir)))
		err = http.Serve(ln, nil)
		if err != nil {
			print(err.Error())
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
		URL:       fmt.Sprintf("http://%s/index.html", ln.Addr().String()),
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
