package unfsd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"gearbox/global"
	"gearbox/help"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"github.com/getlantern/errors"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Unfsd struct {
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	EntityParent   messages.MessageAddress
	Boxname        string
	State          *states.Status
	Task           *tasks.Task

	ExportsBaseDir *ospaths.Dir
	ExportsJson    *ospaths.File
	ExportsFile    *ospaths.File
	PidFile        *ospaths.File
	//NfsCmd         string
	//NfsArgs        []string
	Server         Server
	//State          UnfsdState

	// State polling delays.
	NoWait      bool
	WaitDelay   time.Duration
	WaitRetries int

	mutex           sync.RWMutex // Mutex control for map.
	Channels        *channels.Channels
	//channelHandler  *channels.Subscriber
	OsPaths         *ospaths.BasePaths
}
type Args Unfsd

type ExportData struct {
	MountPoint string
	Options    []string
}

// Server manages exporting an NFS mount.
type Server struct {
	sync.Mutex                           `json:"-"`
	BasePath         string              `json:"basePath"`
	ExportedName     string              `json:"exportedName"`
	ExportedNamePath string              `json:"exportedNamePath"`
	ExportOptions    string              `json:"exportOptions"`
	Network          string              `json:"network"`
	ClientIPs        map[string]struct{} `json:"clients"`
	Volumes          map[string]int32    `json:"volumes"`
	Exported         map[string]struct{} `json:"exported"`
	ClientValidator  NfsClientValidator  `json:"clientValidator"`
}

type UnfsdState struct {
	Name         string
	LastSts      status.Status
	CurrentState int
	WantState    int
}


// Refine facade.DfsValidator to avoid circular dependencies
type NfsClientValidator interface {
	ValidateClient(string) bool
}

var (
	// ErrInvalidExportedName is returned when an exported name is not a valid single directory name
	ErrInvalidExportedName = fmt.Errorf("nfs server: invalid exported name")
	// ErrInvalidBasePath is returned when the local path to export is invalid
	ErrInvalidBasePath = fmt.Errorf("nfs server: invalid base path")
	// ErrBasePathNotDir is returned when the base path is not a directory
	ErrBasePathNotDir = fmt.Errorf("nfs server: base path not a directory")
	// ErrInvalidNetwork is returned when the network specifier does not parse in CIDR format
	ErrInvalidNetwork = fmt.Errorf("nfs server: the network value is not CIDR")
)

var (
	osMkdirAll = os.MkdirAll
	osChmod    = os.Chmod
	fsidIdx    int32
)

const defaultDirectoryPerm = 0755
const etcExportsStartMarker = "# -EXPORTS START-\n"
const etcExportsEndMarker = "# -EXPORTS END-\n"
const etcExportsRemoveComment = "# export removed: "

var ExportTemplate = `{{range .}}
{{.MountPoint}} {{range .Options}} {{.Options}} {{end}}
{{end}}`

const DefaultBaseDir		= "unfsd"
const DefaultExportsFile	= "unfsd/etc/exports"
const DefaultExportsJson	= "unfsd/etc/exports.json"
const DefaultPidFile		= "unfsd/unfsd.pid"


func New(args ...Args) (*Unfsd, status.Status) {

	var sts status.Status
	var err error
	var _args Args
	unfsd := &Unfsd{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Channels == nil {
			err = _args.EntityId.ProduceError("channels pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = _args.EntityId.ProduceError("ospaths is nil")
			break
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.ExportsBaseDir == nil {
			// _args.ExportsBaseDir = &_args.OsPaths.ProjectBaseDir
			_args.ExportsBaseDir = &_args.OsPaths.UserHomeDir
		}
		_, err = _args.ExportsBaseDir.CreateIfNotExists()
		if err != nil {
			break
		}

		if _args.ExportsFile == nil {
			_args.ExportsFile = _args.OsPaths.EventBrokerDir.AddFileToPath(DefaultExportsFile)
			// filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsBridge.GetAdminRootDir(), DefaultExportsFile))
		}

		if _args.ExportsJson == nil {
			_args.ExportsJson = _args.OsPaths.EventBrokerDir.AddFileToPath(DefaultExportsJson)
			// filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsBridge.GetAdminRootDir(), DefaultExportsJson))
		}

		_args.PidFile = _args.OsPaths.EventBrokerDir.AddFileToPath(DefaultPidFile)


		fmt.Printf("ExportsFile:%s\n", _args.ExportsFile)

		*unfsd = Unfsd(_args)

		// Check exports file access.
		err = unfsd.ReadExport()
		if err != nil {
			break
		}

	}

	return unfsd, sts
}


func (me *Unfsd) ReadExport() error {

	// Ensure we read the JSON export file and update the UNFSD exports file.
	// UNFSD exports file will ALWAYS be updated to reflect the JSON file.
	var err error

	for range only.Once {
		err = EnsureNotNil(me)
		if err != nil {
			break
		}

		err = me.readJsonExport()
		if err != nil {
			break
		}

		err = me.writeNfsExport()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Unfsd) WriteExport() error {

	// Ensure we write the JSON export file and update the UNFSD exports file.
	var err error

	for range only.Once {
		err = EnsureNotNil(me)
		if err != nil {
			break
		}

		err = me.writeJsonExport()
		if err != nil {
			break
		}

		err = me.writeNfsExport()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Unfsd) readJsonExport() error {

	//var sts status.Status
	var content []byte
	var err error

	for range only.Once {
		err = EnsureNotNil(me)
		if err != nil {
			break
		}

		content, err = ReadFile(me.ExportsJson)
		if err != nil {
			break
		}
		content = []byte(me.ParsePaths(string(content)))

		if err != nil {
			// Or use me.OsBridge.GetProjectBaseDir()
			me.NewServer(me.ExportsBaseDir, "Sites", "0.0.0.0/0")
			me.Server.AddVolume("default")
			err = me.WriteExport()
			break
		}

		err = json.Unmarshal(content, &me.Server)
		if err != nil {
			//sts = status.Fail(&status.Args{
			//	Message: fmt.Sprintf("UNFSD: Failed to read JSON file '%s' - %v", me.ExportsJson, err),
			//	Help: fmt.Sprintf("Ensure '%s' is in correct format per %s"),
			//	Data: err,
			//})
			break
		}

		if me.Server.BasePath != me.ExportsBaseDir.String() {
			// If we have overridden the BasePath in the args.
			me.Server.BasePath = me.ExportsBaseDir.String()
		}
	}

	return err
}

func (me *Unfsd) writeJsonExport() status.Status {

	var sts status.Status
	var err error
	var content []byte

	for range only.Once {

		content, err = json.Marshal(&me.Server)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Failed to write JSON file '%s' - %v", me.ExportsJson, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		err, sts = WriteFile(me.ExportsJson, []byte(content), 0664)
		if err != nil {
			// WARN
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Could not look up deprecated exports path %s: %s", me.ExportsJson, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return sts
}

func (me *Unfsd) writeNfsExport() status.Status {

	var sts status.Status
	var err error
	var serviced_exports string

	for range only.Once {

		network := me.Server.Network
		if network == "0.0.0.0/0" {
			network = ""	// "*" // turn this in to nfs 'allow all hosts' syntax
		}

		if err := os.MkdirAll(me.ExportsBaseDir.String(), 0775); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", me.ExportsBaseDir, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		//edir := filepath.Join(me.ExportsBaseDir.String(), me.Server.ExportedName)
		edir := me.ExportsBaseDir.AddToPath(me.Server.ExportedName)
		if err = os.MkdirAll(edir.String(), 0775); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", me.ExportsBaseDir, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		exports := make(map[string]struct{})
		serviced_exports = fmt.Sprintf("%s\t%s(rw,fsid=0,no_root_squash,insecure,no_subtree_check,async,crossmnt)\n",
			me.ExportsBaseDir.AddToPath(me.Server.ExportedName).String(), network)
		for volume, fsid := range me.Server.Volumes {
			volume = filepath.Clean(volume)
			_, volName := filepath.Split(volume)
			exports[volName] = struct{}{}
			exported := filepath.Join(edir.String(), volName)

			serviced_exports += fmt.Sprintf("%s\t%s(rw,fsid=%d,no_root_squash,insecure,no_subtree_check,async)\n",
				exported, network, fsid)
		}
		me.Server.Exported = exports

		originalContents, err := readFileIfExists(me.ExportsFile)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error reading exports file '%s': %v", me.ExportsBaseDir, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		// comment out lines that conflicts with serviced exported mountpoints
		mountpaths := map[string]bool{me.ExportsBaseDir.String(): true, filepath.Join(me.ExportsBaseDir.String(), me.Server.ExportedName): true}
		filteredContent := ""
		scanner := bufio.NewScanner(strings.NewReader(originalContents))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "#") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					mountpoint := fields[0]
					if _, ok := mountpaths[mountpoint]; ok {
						filteredContent += etcExportsRemoveComment + line + "\n"
						continue
					}
				}
			}

			filteredContent += line + "\n"
		}
		// fmt.Printf("filteredContent:%v\n", filteredContent)

		// create file content
		preamble, postamble := filteredContent, ""
		if index := strings.Index(filteredContent, etcExportsStartMarker); index >= 0 {
			preamble = filteredContent[:index]
			remainder := filteredContent[index:]
			if index := strings.Index(remainder, etcExportsEndMarker); index >= 0 {
				postamble = remainder[index+len(etcExportsEndMarker):]
			}
		}
		fileContents := preamble + etcExportsStartMarker + serviced_exports + etcExportsEndMarker + postamble

		err, sts = WriteFile(me.ExportsFile, []byte(fileContents), 0664)
	}

	// Force UNFSD to re-read exports file.
	if me.getPid() > 0 {
		//syscall.Kill(me.getPid(), syscall.SIGHUP)
	}

	sts = status.Success("UNFSD: Exported:\n %s", serviced_exports)

	return sts
}


func (me *Unfsd) ParsePaths(i string) string {

	if me.OsPaths == nil {
		return i
	}

	strReplace := map[string]string {
		"{{.LocalDir}}":              me.OsPaths.LocalDir.String(),
		"{{.UserHomeDir}}":           me.OsPaths.UserHomeDir.String(),
		"{{.AdminRootDir}}":          me.OsPaths.AdminRootDir.String(),
		"{{.CacheDir}}":              me.OsPaths.CacheDir.String(),
		"{{.ProjectBaseDir}}":        me.OsPaths.ProjectBaseDir.String(),
		"{{.UserConfigDir}}":         me.OsPaths.UserConfigDir.String(),
		"{{.EventBrokerDir}}":        me.OsPaths.EventBrokerDir.String(),
		"{{.EventBrokerWorkingDir}}": me.OsPaths.EventBrokerWorkingDir.String(),
		"{{.EventBrokerLogDir}}":     me.OsPaths.EventBrokerLogDir.String(),
		"{{.EventBrokerEtcDir}}":     me.OsPaths.EventBrokerEtcDir.String(),
		"{{.Platform}}":              runtime.GOOS + "_" + runtime.GOARCH,
	}

	for k, v := range strReplace {
		i = strings.ReplaceAll(i, k, v)
	}

	return i
}

func (me *Unfsd) readNfsExport() status.Status {

	// Stub method.
	// We won't ever need to read the UNFSD exports file because it will ALWAYS
	// be a translated copy of the exports.json file.

	return nil
}


////////////////////////////////////////////////////////////////////////////////

// NewServer returns a Unfsd.Server object that manages the given nfs mounts to
// configured clients;  basePath is the path for volumes, exportedName is the container dir to hold exported volumes
func (me *Unfsd) NewServer(basePath *ospaths.Dir, exportedName string, network string) error {

	if len(exportedName) < 2 || strings.Contains(exportedName, "/") {
		return ErrInvalidExportedName
	}

	if len(basePath.String()) < 2 {
		return ErrInvalidBasePath
	}

	if err := verifyExportsBaseDir(basePath); err != nil {
		return err
	}

	// exportedNamePath := filepath.Join(me.ExportsBaseDir.String(), exportedName)
	ename := me.ExportsBaseDir.AddToPath(exportedName)
	if err := verifyExportsBaseDir(ename); err != nil {
		return err
	}

	if _, _, err := net.ParseCIDR(network); err != nil {
		return ErrInvalidNetwork
	}

	me.Server = Server{
		BasePath:         basePath.String(),
		ExportedName:     exportedName,
		ExportedNamePath: ename.String(),
		ExportOptions:    "rw,insecure,no_subtree_check,async",
		ClientIPs:        make(map[string]struct{}),
		Network:          network,
		Volumes:          make(map[string]int32),
		Exported:         make(map[string]struct{}),
		ClientValidator:  nil,
	}

	return nil
}


// Reload ensures that the nfs exports are visible to all clients
func (me *Unfsd) Reload() status.Status {

	var sts status.Status

	me.Server.Lock()
	defer me.Server.Unlock()

	// Not implemented yet.

	return sts
}


// Restart restarts the nfs subsystem
func (me *Unfsd) Restart() status.Status {

	var sts status.Status

	for range only.Once {

		sts = me.writeNfsExport()
		if is.Error(sts) {
			break
		}

		sts = me.Stop()
		if is.Error(sts) {
			break
		}

		sts = me.Start()
		if is.Error(sts) {
			break
		}
	}

	return sts
}


// Stop stops the nfs subsystem
func (me *Unfsd) Stop() status.Status {

	var sts status.Status

	//me.Server.Lock()
	//defer me.Server.Unlock()
	//
	//for range only.Once {
	//
	//	if !me.Daemon.IsLoaded() {
	//		sts = status.Success("UNFSD: Stopped.")
	//		break
	//	}
	//
	//	me.State.WantState = external.StatePowerOff
	//	if err := me.Daemon.Unload(); err != nil {
	//		sts = status.Fail(&status.Args{
	//			Message: fmt.Sprintf("UNFSD: Error stopping: %v", err),
	//			Help:    help.ContactSupportHelp(), // @TODO need better support here
	//			Data:    err,
	//		})
	//		break
	//	}
	//
	//	me.State, sts = me.GetState()
	//}

	return sts
}


// Stop stops the nfs subsystem
func (me *Unfsd) Start() status.Status {

	var sts status.Status
	//var err error
	//
	//me.Server.Lock()
	//defer me.Server.Unlock()
	//
	//for range only.Once {
	//
	//	if me.Daemon.IsLoaded() {
	//		sts = status.Success("UNFSD: Started.")
	//		break
	//	}
	//
	//	sts = me.writeNfsExport()
	//	if is.Error(sts) {
	//		break
	//	}
	//
	//	me.State.WantState = external.StateRunning
	//	if err = me.Daemon.Load(); err != nil {
	//		sts = status.Fail(&status.Args{
	//			Message: fmt.Sprintf("UNFSD: Error starting: %v", err),
	//			Help:    help.ContactSupportHelp(), // @TODO need better support here
	//			Data:    err,
	//		})
	//		break
	//	}
	//
	//	me.State, sts = me.GetState()
	//}

	return sts
}


//func (me *Unfsd) GetState() (state UnfsdState, sts status.Status) {
//
//	for range only.Once {
//		sts = EnsureNotNil(me)
//		if is.Error(sts) {
//			break
//		}
//
//		if me.State.WantState == external.StateInit {
//			me.State.WantState = me.State.CurrentState
//		}
//
//		if !me.Daemon.IsLoaded() {
//			me.State.LastSts = status.Warn("%s UNFSD - not loaded", global.Brandname)
//			me.State.CurrentState = external.StateLoaded
//			break
//		}
//
//		if !me.Daemon.IsRunning() {
//			me.State.LastSts = status.Warn("%s UNFSD - halted", global.Brandname)
//			me.State.CurrentState = external.StatePowerOff
//			break
//		}
//
//		me.State.LastSts = status.Success("%s UNFSD - running", global.Brandname)
//		me.State.CurrentState = external.StateRunning
//	}
//
//	return me.State, me.State.LastSts
//}


func (me *Unfsd) getPid() (pid int) {

	var exists bool
	var data []byte
	var err error

	if exists, err = doesExist(me.PidFile); err != nil {
		return 0
	}

	if exists {
		data, err = ioutil.ReadFile(me.PidFile.String())
		if err != nil {
			return 0
		}

		pid, err = strconv.Atoi(string(data))
		if err != nil {
			return 0
		}
	}

	return pid
}



////////////////////////////////////////////////////////////////////////////////
// Server methods

// ExportPath returns the external export name; foo for nfs export /exports/foo
func (me *Server) ExportPath() string {

	return filepath.Join("/", me.ExportedName)
}


// Returns the export path name; a combination of the me.ExportsBaseDir and ExportPath
func (me *Server) ExportNamePath() string {

	return me.ExportedNamePath
}


// Clients returns the IP Addresses of the current clients
func (me *Server) Clients() []string {

	clients := make([]string, len(me.ClientIPs))

	i := 0
	for key := range me.ClientIPs {
		clients[i] = key
	}

	return clients
}


func (me *Server) SetClientValidator(validator NfsClientValidator) {

	me.ClientValidator = validator
}


// SetClients replaces the existing clients with the new clients
func (me *Server) SetClients(clients ...string) {

	me.Lock()
	defer me.Unlock()

	filteredClients := me.filterHostsWithoutPerms(clients)
	me.ClientIPs = make(map[string]struct{})

	for _, client := range filteredClients {
		me.ClientIPs[client] = struct{}{}
	}

}


// VolumeCreated set that path of a volume that should be exported
func (me *Server) AddVolume(volumePath string) error {

	me.Lock()
	defer me.Unlock()

	fsid := atomic.AddInt32(&fsidIdx, 1)
	me.Volumes[volumePath] = fsid

	return nil
}


// VolumeCreated set that path of a volume that should be exported
func (me *Server) RemoveVolume(volumePath string) error {

	me.Lock()
	defer me.Unlock()

	delete(me.Volumes, volumePath)

	return nil
}


func (me *Server) filterHostsWithoutPerms(clients []string) []string {

	if me.ClientValidator == nil {
		return clients
	}

	filteredClients := []string{}

	for _, client := range clients {
		if me.ClientValidator.ValidateClient(client) {
			filteredClients = append(filteredClients, client)
		} else {
			// fmt.Printf("Filtered NFS client with ip %s", client)
		}
	}

	return filteredClients
}



////////////////////////////////////////////////////////////////////////////////
// Misc functions.

func EnsureNotNil(bx *Unfsd) error {
	var err error

	if bx == nil {
		err = errors.New("nil structure")
	}

	return err
}


func verifyExportsBaseDir(path *ospaths.Dir) error {
	stat, err := os.Stat(path.String())
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//handle does not exist
		return osMkdirAll(path.String(), defaultDirectoryPerm)
	}
	if !stat.IsDir() {
		return ErrBasePathNotDir
	}
	if (stat.Mode() & defaultDirectoryPerm) != defaultDirectoryPerm {
		err = osChmod(path.String(), defaultDirectoryPerm)
	}
	return err
}


func readFileIfExists(path *ospaths.File) (s string, err error) {

	var exists bool

	if exists, err = doesExist(path); err != nil {
		return s, err
	}

	if exists {
		bytes, err := ioutil.ReadFile(path.String())
		if err != nil {
			return s, err
		}
		s = string(bytes)
	}

	return s, nil
}


// removeDeprecated will clean up any old exports path
//func removeDeprecated(dirpath *ospaths.Dir) status.Status {
//	var sts status.Status
//
//	for range only.Once {
//
//		if dirContents, err := ioutil.ReadDir(dirpath); !os.IsNotExist(err) {
//			if err != nil {
//				// WARN
//				sts = status.Fail(&status.Args{
//					Message: fmt.Sprintf("UNFSD: Could not look up deprecated exports path %s: %s", dirpath, err),
//					Help:    help.ContactSupportHelp(), // @TODO need better support here
//					Data:    err,
//				})
//				break
//			}
//
//			// Guarantee the path is now empty
//			if l := len(dirContents); l > 0 {
//				// WARN
//				sts = status.Fail(&status.Args{
//					Message: fmt.Sprintf("UNFSD: Path %s is not empty.", dirpath),
//					Help:    help.ContactSupportHelp(), // @TODO need better support here
//					Data:    err,
//				})
//				break
//
//			} else {
//				// Remove the path only if it is empty
//				if err := os.Remove(dirpath); err != nil {
//					// WARN
//					sts = status.Fail(&status.Args{
//						Message: fmt.Sprintf("UNFSD: Could not remove deprecated path %s: %s", dirpath, err),
//						Help:    help.ContactSupportHelp(), // @TODO need better support here
//						Data:    err,
//					})
//					break
//				}
//
//				sts = status.Success(fmt.Sprintf("UNFSD: Deleted deprecated path %s", dirpath))
//			}
//		}
//	}
//
//	return sts
//}


func doesExist(path *ospaths.File) (bool, error) {

	_, err := os.Stat(path.String())
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}


func ReadFile(filename *ospaths.File) ([]byte, error) {

	var data []byte
	var err error

	for range only.Once {

		err = filename.FileExists()
		if err != nil {
			err = errors.New(fmt.Sprintf("UNFSD: Missing file '%s'", filename))
			break
		}

		data, err = ioutil.ReadFile(filename.String())
		if err != nil {
			err = errors.New(fmt.Sprintf("UNFSD: Error reading file '%s'", filename))
			break
		}
	}

	return data, err
}


// WriteFile will write the given data to the filename in an atomic manner so that
// partial writes are not possible.
func WriteFile(filename *ospaths.File, data []byte, perm os.FileMode) (error, status.Status) {

	var sts status.Status
	var err error

	for range only.Once {

		// find the dirname of the filename
		d := filepath.Dir(filename.String())
		if err = os.MkdirAll(d, 0755); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", d, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		tempfile, err := ioutil.TempFile(d, filepath.Base(filename.String()))
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating temp file '%s': %v", tempfile, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
		name := tempfile.Name()
		defer os.Remove(name)

		if err = tempfile.Close(); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error closing tempfile '%s': %v", name, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		if err = ioutil.WriteFile(name, data, perm); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error writing file '%s': %v", name, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		if err = os.Chmod(name, perm); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error chmoding file '%s': %v", name, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		if err = os.Rename(name, filename.String()); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error renaming file '%s' ->: %v", name, filename, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return err, sts
}

