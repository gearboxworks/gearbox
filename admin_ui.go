package gearbox

import (
	"context"
	"fmt"
	"gearbox/host"
	"gearbox/only"
	"github.com/labstack/gommon/log"
	"github.com/zserge/lorca"
	"github.com/zserge/webview"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

//
// [1] https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
// [2] https://github.com/zserge/webview
// [3] https://github.com/zserge/lorca
//

const (
	UiWindowTitle = "Gearbox"
	UiHeight      = 600
	UiWidth       = 800
	UiResizable   = true
	UiDebug       = true
)

type UiWindow struct {
	Title        string
	Height       int
	Width        int
	NotResizable bool
	NoDebug      bool
}

type ViewerType string

const (
	DefaultViewer ViewerType = LorcaViewer
	WebViewViewer ViewerType = "webview"
	LorcaViewer   ViewerType = "lorca"
)

type AdminUi struct {
	ViewerType    ViewerType
	webListener   net.Listener
	HostConnector host.Connector
	Gearbox       *Gearbox
	webServer     *http.Server
	api           *HostApi
	Window        *UiWindow
}

type UiWindowArgs struct {
	Title  string
	Height int
	Width  int
}

func NewUiWindow(args *UiWindowArgs) *UiWindow {
	if args.Title != "" {
		// If args.Title contains %s it will replace with UiWindowTitle's value
		// If not, it will just use value of args.Title
		args.Title = fmt.Sprintf(args.Title, UiWindowTitle)
	}
	if args.Height == 0 {
		args.Height = UiHeight
	}
	if args.Width == 0 {
		args.Width = UiWidth
	}
	return &UiWindow{
		Title:  args.Title,
		Height: args.Height,
		Width:  args.Width,
	}
}

func NewAdminUi(gearbox *Gearbox, viewer ViewerType) *AdminUi {
	ui := AdminUi{
		Gearbox:       gearbox,
		HostConnector: gearbox.HostConnector,
		ViewerType:    viewer,
		Window: NewUiWindow(&UiWindowArgs{
			Title: "%s - " + fmt.Sprintf("[%s]", viewer),
		}),
	}
	return &ui
}

func (me *AdminUi) Initialize() {
	me.webListener = me.GetWebListener()
	me.webServer = me.GetWebServer()
	me.api = NewHostApi(me.Gearbox)
	me.WriteAssetsToAdminWebRoot()
}

func (me *AdminUi) WriteAssetsToAdminWebRoot() {
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
	me.WriteApiBaseUrls()
}

func (me *AdminUi) StartApi() {
	go me.api.Start()
}

func (me *AdminUi) Start() {
	go me.StartApi()
	go me.ServeWeb()
	switch me.ViewerType {
	case WebViewViewer:
		me.StartWebView()
	case LorcaViewer:
		me.StartLorca()
	default:
		log.Warnf("invalid viewer type '%s'.", me.ViewerType)
	}
}

func (me *AdminUi) StartLorca() {
	win := me.Window
	//time.Sleep(time.Second*3)
	ui, err := lorca.New(
		me.GetWebRootFileUrl(),
		string(me.GetWebRootDir()),
		win.Width,
		win.Height,
	)
	if err != nil {
		log.Warnf("error loading Lorca to view Gearbox Admin UI: %s", err)
	}
	<-ui.Done()
}

func (me *AdminUi) StartWebView() {
	win := me.Window
	wv := webview.New(webview.Settings{
		Title:     win.Title,
		Height:    win.Height,
		Width:     win.Width,
		Resizable: !win.NotResizable,
		Debug:     !win.NoDebug,
		URL:       me.GetWebRootFileUrl(),
	})
	wv.Run()
}

func (me *AdminUi) WriteApiBaseUrls() {
	var err error
	url := me.api.Url()
	file := me.GetApiBaseUrls()
	err = ioutil.WriteFile(file, NewApiBaseUrls(url, url).Bytes(), os.ModePerm)
	if err != nil {
		log.Warnf("error writing API bootrap file '%s': %s",
			me.GetApiBaseUrls(),
			err,
		)
	}
}

func (me *AdminUi) GetApiBaseUrls() string {
	return fmt.Sprintf("%s/api.json", me.GetWebRootDir())
}

func (me *AdminUi) GetWebListener() net.Listener {
	var err error
	for range only.Once {
		if me.webListener != nil {
			break
		}
		me.webListener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Warnf("error initiating a TCP connection for AdminUi on '127.0.0.0:0': %s", err)
		}
		if me.Gearbox.Options.IsDebug {
			fmt.Printf("\nListening on %s", me.GetHostname())
		}
	}
	return me.webListener
}

func (me *AdminUi) GetHostname() string {
	return me.webListener.Addr().String()
}

func (me *AdminUi) GetWebRootUrl() string {
	return fmt.Sprintf("http://%s", me.GetHostname())
}

func (me *AdminUi) GetWebRootFileUrl() string {
	return fmt.Sprintf("%s/index.html", me.GetWebRootUrl())
}

func (me *AdminUi) GetWebRootDir() http.Dir {
	return http.Dir(me.HostConnector.GetAdminRootDir())
}

func (me *AdminUi) GetWebRootFileDir() string {
	return fmt.Sprintf("%s/index.html", me.GetWebRootDir())
}

func (me *AdminUi) GetWebHandler() http.Handler {
	return addCorsMiddleware(me.GetWebRootUrl(), http.FileServer(me.GetWebRootDir()))
}

func addCorsMiddleware(rooturl string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", rooturl)
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
		next.ServeHTTP(w, r)
	})
}

func shutdownServer(srv *http.Server) {
	fmt.Print("Shutting down Gearbox Admin Console web server...\n")
	err := srv.Shutdown(context.TODO())
	if err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func (me *AdminUi) GetWebServer() *http.Server {
	for range only.Once {
		if me.webServer != nil {
			break
		}
		addr := me.GetHostname()
		me.webServer = &http.Server{
			Addr:    addr,
			Handler: me.GetWebHandler(),
		}
	}
	return me.webServer
}

func (me *AdminUi) GetHostApi() *HostApi {
	for range only.Once {
		if me.api != nil {
			break
		}
		me.api = NewHostApi(me.Gearbox)
	}
	return me.api
}

func (me *AdminUi) Close() {
	shutdownServer(me.webServer)
	//err := me.webListener.Close()
	//if err != nil {
	//	log.Warnf("error attempting to close AdminUi: %s", err)
	//}
	me.api.Stop()
}

func (me *AdminUi) ServeWeb() {
	err := me.webServer.Serve(me.webListener)
	// returns ErrServerClosed on graceful close
	if err != http.ErrServerClosed {
		// NOTE: there is a chance that next line won't have time to run,
		// as main() doesn't wait for this goroutine to stop. don't use
		// code with race conditions like these for production. see post
		// comments below on more discussion on how to handle this.
		log.Fatalf("error closing http.Server in AdminUi.ServeWeb(): %s", err)
	}
}
