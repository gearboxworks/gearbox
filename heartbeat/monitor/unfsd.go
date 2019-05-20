package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gearbox/global"
	"gearbox/heartbeat/daemon"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)


type Unfsd struct {
	Boxname        string
	ExportsBaseDir string
	ExportsJson    string
	ExportsFile    string
	PidFile        string
	NfsCmd         string
	NfsArgs        []string
	Server         Server
	Daemon         *daemon.Daemon

	// State polling delays.
	NoWait      bool
	WaitDelay   time.Duration
	WaitRetries int

	OsSupport oss.OsSupporter
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

const DefaultExportsFile	= "heartbeat/unfsd/etc/exports"
const DefaultExportsJson	= "heartbeat/unfsd/etc/exports.json"
const DefaultNfsBin			= "heartbeat/unfsd/bin/unfsd"
const DefaultPidFile		= "heartbeat/unfsd/etc/unfsd.pid"


func NewUnfsd(OsSupport oss.OsSupporter, args ...Args) (*Unfsd, status.Status) {

	var sts status.Status
	var _args Args
	unfsd := &Unfsd{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.ExportsBaseDir == "" {
			_args.ExportsBaseDir = string(OsSupport.GetUserHomeDir())
		}

		if _args.ExportsFile == "" {
			_args.ExportsFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), DefaultExportsFile))
		}

		if _args.ExportsJson == "" {
			_args.ExportsJson = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), DefaultExportsJson))
		}

		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), DefaultPidFile))
		_args.NfsCmd = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), DefaultNfsBin))

		// Need to check for existance of UNFSD pid.

		_args.NfsArgs = []string {
			"-e",
			_args.ExportsFile,

			"-i",
			_args.PidFile,

			// "-u",		// Use unprivileged port for NFS service - May be required for non-admin access.

			// "-n",		// Use specified port for NFS service.
			// "2049",		// Use specified port for NFS service.

			// "-m",		// Use specified port for MOUNT service.
			// "2049",		// Use specified port for MOUNT service.

			// "-t",		// Use TCP instead of UDP.

			// "-p",		// Don't register with the portmapper - good for security.

			// "-c",		// Enable clustering extension - good for minor changes across files.

			// "-C",		// Enable clustering extension - good for minor changes across files.
			// "path",		// Enable clustering extension - good for minor changes across files.

			"-s",			// Single user mode - force all UID/GID mapping to be that of the client side.

			// "-b",		// Brute force file searching - Makes NFS really slow, but avoids .nfsd* files being created.

			"-l",			// Bind to specific address.
			"0.0.0.0",		// Bind to specific address.

			"-d",			// Don't detach from terminal.

			// "-r",		// Report unreadable executables as readable - usually not required to be enabled.

			// "-T",		// Just test exports file for issues.
		}

		fmt.Printf("ExportsFile:%s\n", _args.ExportsFile)
		fmt.Printf("NfsCmd:%s\n", _args.NfsCmd)

		// Start a new UNFSD instance.
		_args.Daemon = daemon.NewDaemon(_args.OsSupport, daemon.Args{
			Boxname: _args.Boxname,
			ServiceData: daemon.PlistData {
				Label: "com.gearbox.unfsd",
				Program:   _args.NfsCmd,
				ProgramArgs: _args.NfsArgs,
				KeepAlive: true,
				RunAtLoad: true,
			},
		})

		*unfsd = Unfsd(_args)

		// Check exports file access.
		sts = unfsd.ReadExport()
		if is.Error(sts) {
			break
		}

		// Check nfsd binary access.
		_, sts = ReadFile(_args.NfsCmd)
		if is.Error(sts) {
			break
		}

	}

	return unfsd, sts
}


func (me *Unfsd) ReadExport() (status.Status) {

	// Ensure we read the JSON export file and update the UNFSD exports file.
	// UNFSD exports file will ALWAYS be updated to reflect the JSON file.
	var sts status.Status

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.readJsonExport()
		if is.Error(sts) {
			break
		}

		sts = me.writeNfsExport()
		if is.Error(sts) {
			break
		}
	}

	return sts
}


func (me *Unfsd) WriteExport() (status.Status) {

	// Ensure we write the JSON export file and update the UNFSD exports file.
	var sts status.Status

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.writeJsonExport()
		if is.Error(sts) {
			break
		}

		sts = me.writeNfsExport()
		if is.Error(sts) {
			break
		}
	}

	return sts
}


func (me *Unfsd) readJsonExport() (status.Status) {

	var sts status.Status
	var content []byte
	var err error

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		content, sts = ReadFile(me.ExportsJson)
		if is.Error(sts) {
			break
		}

		if is.Warn(sts) {
			me.NewServer(me.ExportsBaseDir, "Sites", "0.0.0.0/0")
			me.Server.AddVolume("default")
			sts = me.WriteExport()
			break
		}

		err = json.Unmarshal(content, &me.Server)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Failed to read JSON file '%s' - %v", me.ExportsJson, err),
				Help: fmt.Sprintf("Ensure '%s' is in correct format per %s"),
				Data: err,
			})
			break
		}

		if me.Server.BasePath != me.ExportsBaseDir {
			// If we have overridden the BasePath in the args.
			me.Server.BasePath = me.ExportsBaseDir
		}
	}

	return sts
}


func (me *Unfsd) writeJsonExport() (status.Status) {

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


func (me *Unfsd) writeNfsExport() (status.Status) {

	var sts status.Status
	var err error
	var serviced_exports string

	for range only.Once {

		network := me.Server.Network
		if network == "0.0.0.0/0" {
			network = ""	// "*" // turn this in to nfs 'allow all hosts' syntax
		}

		if err := os.MkdirAll(me.ExportsBaseDir, 0775); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", me.ExportsBaseDir, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		edir := filepath.Join(me.ExportsBaseDir, me.Server.ExportedName)
		if err = os.MkdirAll(edir, 0775); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", me.ExportsBaseDir, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		exports := make(map[string]struct{})
		serviced_exports = fmt.Sprintf("%s\t%s(rw,fsid=0,no_root_squash,insecure,no_subtree_check,async,crossmnt)\n",
			me.ExportsBaseDir + "/" + me.Server.ExportedName, network)
		for volume, fsid := range me.Server.Volumes {
			volume = filepath.Clean(volume)
			_, volName := filepath.Split(volume)
			exports[volName] = struct{}{}
			exported := filepath.Join(edir, volName)

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
		mountpaths := map[string]bool{me.ExportsBaseDir: true, filepath.Join(me.ExportsBaseDir, me.Server.ExportedName): true}
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
		fmt.Printf("filteredContent:%v\n", filteredContent)

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
		syscall.Kill(me.getPid(), syscall.SIGHUP)
	}

	sts = status.Success("UNFSD: Exported:\n %s", serviced_exports)

	return sts
}


func (me *Unfsd) readNfsExport() (status.Status) {

	// Stub method.
	// We won't ever need to read the UNFSD exports file because it will ALWAYS
	// be a translated copy of the exports.json file.

	return nil
}


////////////////////////////////////////////////////////////////////////////////

// NewServer returns a Unfsd.Server object that manages the given nfs mounts to
// configured clients;  basePath is the path for volumes, exportedName is the container dir to hold exported volumes
func (me *Unfsd) NewServer(basePath, exportedName, network string) (error) {

	if len(exportedName) < 2 || strings.Contains(exportedName, "/") {
		return ErrInvalidExportedName
	}

	if len(basePath) < 2 {
		return ErrInvalidBasePath
	}

	if err := verifyExportsBaseDir(basePath); err != nil {
		return err
	}

	exportedNamePath := filepath.Join(me.ExportsBaseDir, exportedName)
	if err := verifyExportsBaseDir(exportedNamePath); err != nil {
		return err
	}

	if _, _, err := net.ParseCIDR(network); err != nil {
		return ErrInvalidNetwork
	}

	me.Server = Server{
		BasePath:         basePath,
		ExportedName:     exportedName,
		ExportedNamePath: exportedNamePath,
		ExportOptions:    "rw,insecure,no_subtree_check,async",
		ClientIPs:        make(map[string]struct{}),
		Network:          network,
		Volumes:          make(map[string]int32),
		Exported:         make(map[string]struct{}),
		ClientValidator:  nil,
	}

	return nil
}


// Sync ensures that the nfs exports are visible to all clients
func (me *Unfsd) Sync() (error, status.Status) {

	var sts status.Status
	var err error

	me.Server.Lock()
	defer me.Server.Unlock()

	for range only.Once {

		sts = me.writeNfsExport()
		if is.Error(sts) {
			break
		}

		if err = start(); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error starting: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		if err = reload(); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error reloading: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return nil, sts
}


// Restart restarts the nfs subsystem
func (me *Unfsd) Restart() (error, status.Status) {

	var sts status.Status
	var err error

	me.Server.Lock()
	defer me.Server.Unlock()

	for range only.Once {

		sts = me.writeNfsExport()
		if is.Error(sts) {
			break
		}

		if err = restart(); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error restarting: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return nil, sts
}


// Stop stops the nfs subsystem
func (me *Unfsd) Stop() (error, status.Status) {

	var sts status.Status

	me.Server.Lock()
	defer me.Server.Unlock()

	for range only.Once {

		if err := stop(); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error stopping: %v", err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return nil, sts
}


func (me *Unfsd) getPid() (pid int) {

	var exists bool
	var data []byte
	var err error

	if exists, err = doesExist(me.PidFile); err != nil {
		return 0
	}

	if exists {
		data, err = ioutil.ReadFile(me.PidFile)
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

func EnsureNotNil(bx *Unfsd) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    nil,
		})
	}
	return sts
}


func verifyExportsBaseDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//handle does not exist
		return osMkdirAll(path, defaultDirectoryPerm)
	}
	if !stat.IsDir() {
		return ErrBasePathNotDir
	}
	if (stat.Mode() & defaultDirectoryPerm) != defaultDirectoryPerm {
		err = osChmod(path, defaultDirectoryPerm)
	}
	return err
}


func readFileIfExists(path string) (s string, err error) {

	var exists bool

	if exists, err = doesExist(path); err != nil {
		return s, err
	}

	if exists {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return s, err
		}
		s = string(bytes)
	}

	return s, nil
}


// removeDeprecated will clean up any old exports path
func removeDeprecated(dirpath string) status.Status {
	var sts status.Status

	for range only.Once {

		if dirContents, err := ioutil.ReadDir(dirpath); !os.IsNotExist(err) {
			if err != nil {
				// WARN
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("UNFSD: Could not look up deprecated exports path %s: %s", dirpath, err),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    err,
				})
				break
			}

			// Guarantee the path is now empty
			if l := len(dirContents); l > 0 {
				// WARN
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("UNFSD: Path %s is not empty.", dirpath),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    err,
				})
				break

			} else {
				// Remove the path only if it is empty
				if err := os.Remove(dirpath); err != nil {
					// WARN
					sts = status.Fail(&status.Args{
						Message: fmt.Sprintf("UNFSD: Could not remove deprecated path %s: %s", dirpath, err),
						Help:    help.ContactSupportHelp(), // @TODO need better support here
						Data:    err,
					})
					break
				}

				sts = status.Success(fmt.Sprintf("UNFSD: Deleted deprecated path %s", dirpath))
			}
		}
	}

	return sts
}


func doesExist(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}


func ReadFile(filename string) ([]byte, status.Status) {

	var data []byte
	var sts status.Status
	var err error

	for range only.Once {

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			sts = status.Warn("UNFSD: Missing file '%s'.", filename)
			break
		}

		data, err = ioutil.ReadFile(filename)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error reading file '%s'.", filename),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
	}

	return data, sts
}


// WriteFile will write the given data to the filename in an atomic manner so that
// partial writes are not possible.
func WriteFile(filename string, data []byte, perm os.FileMode) (error, status.Status) {

	var sts status.Status
	var err error

	for range only.Once {

		// find the dirname of the filename
		d := filepath.Dir(filename)
		if err = os.MkdirAll(d, 0755); err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("UNFSD: Error creating directory '%s': %v", d, err),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		tempfile, err := ioutil.TempFile(d, filepath.Base(filename))
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

		if err = os.Rename(name, filename); err != nil {
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




////////////////////////////////////////////////////////////////////////////////
// Daemon methods.

func (me *Unfsd) DaemonLoop() (error) {

	return nil
}


func start() (error) {

	return nil
}


func stop() (error) {

	return nil
}


func restart() (error) {

	return nil
}


func reload() (error) {

	return nil
}

