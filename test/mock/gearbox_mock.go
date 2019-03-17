package mock

import (
	"fmt"
	"gearbox"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/host"
	"gearbox/stat"
)

var _ gearbox.Gearbox = (*GearboxObj)(nil)

type GearboxObj struct {
	HostConnector host.Connector
	Config        gearbox.Config
}

type GearboxArgs GearboxObj

func NewGearbox(args *GearboxArgs) gearbox.Gearbox {
	mgb := GearboxObj{}
	mgb = GearboxObj(*args)
	return &mgb
}

func (me *GearboxObj) GetConfig() gearbox.Config {
	return me.Config
}

func (me *GearboxObj) SetConfig(config gearbox.Config) {
	me.Config = config
}

func (me *GearboxObj) GetHostConnector() host.Connector {
	return me.HostConnector
}

func (me *GearboxObj) GetResourceName() api.ResourceName {
	return "a-test-resource"
}

func (me *GearboxObj) GetProjectFilepath(path string, basedir string) (pfp string, status stat.Status) {
	pd, status := me.GetProjectDir(path, basedir)
	return fmt.Sprintf("%s/%s", pd, gearbox.ProjectFilename), status
}

func (me *GearboxObj) GetProjectDir(path string, basedir string) (bd string, status stat.Status) {
	uhd := me.HostConnector.GetUserHomeDir()
	return fmt.Sprintf("%s/%s", uhd, path), status
}

func (me *GearboxObj) ValidateBasedirNickname(nn string, args *gearbox.ValidateArgs) stat.Status {
	return stat.NewOkStatus("nickname '%s' validates just fine during testing", nn)
}

func (me *GearboxObj) StartBox(gearbox.BoxArgs) error {
	panic("implement me")
}

func (me *GearboxObj) StopBox(gearbox.BoxArgs) error {
	panic("implement me")
}

func (me *GearboxObj) PrintBoxStatus(gearbox.BoxArgs) (string, error) {
	panic("implement me")
}

func (me *GearboxObj) RestartBox(gearbox.BoxArgs) error {
	panic("implement me")
}

func (me *GearboxObj) CreateBox(gearbox.BoxArgs) (string, error) {
	panic("implement me")
}

func (me *GearboxObj) ConnectSSH(gearbox.SSHArgs) error {
	panic("implement me")
}

func (me *GearboxObj) Admin(gearbox.ViewerType) {
	panic("implement me")
}

func (me *GearboxObj) Initialize() stat.Status {
	panic("implement me")
}

func (me *GearboxObj) GetStackMap() (gearbox.StackMap, stat.Status) {
	panic("implement me")
}

func (me *GearboxObj) GetGlobalOptions() *gearbox.GlobalOptions {
	panic("implement me")
}

func (me *GearboxObj) GetHostApi() *gearbox.HostApi {
	panic("implement me")
}

func (me *GearboxObj) SetResourceName(api.ResourceName) {
	panic("implement me")
}

func (me *GearboxObj) IsDebug() bool {
	panic("implement me")
}

func (me *GearboxObj) NoCache() bool {
	panic("implement me")
}

func (me *GearboxObj) ProjectExists(string) bool {
	panic("implement me")
}

func (me *GearboxObj) AddBasedir(string, ...string) stat.Status {
	panic("implement me")
}

func (me *GearboxObj) UpdateBasedir(string, string) stat.Status {
	panic("implement me")
}

func (me *GearboxObj) DeleteNamedBasedir(string) stat.Status {
	panic("implement me")
}

func (me *GearboxObj) NamedBasedirExists(string) bool {
	panic("implement me")
}

func (me *GearboxObj) FindProjectWithDetails(string) (*gearbox.Project, stat.Status) {
	panic("implement me")
}

func (me *GearboxObj) AddNamedStackToProject(gearbox.StackName, string) stat.Status {
	panic("implement me")
}

func (me *GearboxObj) RequestAvailableContainers(...*dockerhub.ContainerQuery) (dockerhub.ContainerNames, stat.Status) {
	panic("implement me")
}

func (me *GearboxObj) GetApiUrl(api.ResourceName, api.UriTemplateVars) (url string, status stat.Status) {
	panic("implement me")
}

func (me *GearboxObj) WriteLog([]byte) (int, error) {
	panic("implement me")
}
