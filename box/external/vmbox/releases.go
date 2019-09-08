package vmbox

import (
	"errors"
	"fmt"
	"gearbox/box/external/hypervisor"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"gearbox/global"
	"github.com/cavaliercoder/grab"
	"github.com/gearboxworks/go-status/only"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"os"
	"strings"
	"time"
)

type ReleasesMap map[Version]*Release

func (me *ReleasesMap) EnsureNotNil() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("release is nil")
			break
		}
	}

	return err
}

// @TODO Convention is that plural is slice of singular, e.g. Foos []*Foo
//       Can rename this something else?
type Releases struct {
	Map      ReleasesMap
	Latest   *Release
	Selected *Release
	BaseDir  osdirs.Dir

	channels *channels.Channels
}

func NewReleases(c *channels.Channels) (*Releases, error) {

	var ret *Releases
	var err error

	for range only.Once {
		p := osdirs.New()

		me := Releases{}
		me.BaseDir = p.AppendToUserConfigDir("iso")
		me.Map = make(ReleasesMap)
		me.channels = c

		err = me.UpdateReleases()

		ret = &me

		eblog.Debug(entity.VmBoxEntityName, "created new release structre")
	}

	eblog.LogIfNil(ret, err)
	eblog.LogIfError(err)

	return ret, err
}

func (me *Releases) EnsureNotNil() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("releases is nil")
			break
		}
	}

	return err
}

func (me *Releases) ShowReleases() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		fmt.Printf("Latest release: %v\n\n", me.Latest)
		for _, release := range me.Map {
			fmt.Printf("Assets for release:	%v\n", release.Instance.GetName())
			fmt.Printf("UploadURL: 			%v\n", release.Instance.GetUploadURL())
			fmt.Printf("ZipballURL: 			%v\n", release.Instance.GetZipballURL())
			fmt.Printf("TarballURL: 			%v\n", release.Instance.GetTarballURL())
			fmt.Printf("Body: 				%v\n", release.Instance.GetBody())
			fmt.Printf("AssetsURL: 			%v\n", release.Instance.GetAssetsURL())
			fmt.Printf("URL: 				%v\n", release.Instance.GetURL())
			fmt.Printf("HTMLURL:				%v\n", release.Instance.GetHTMLURL())

			for _, asset := range release.Instance.Assets {
				fmt.Printf("	Name:				%v\n", asset.GetName())
				fmt.Printf("	ID:					%v\n", asset.GetID())
				fmt.Printf("	URL:					%v\n", asset.GetURL())
				fmt.Printf("	Size:				%v\n", asset.GetSize())
				fmt.Printf("	CreatedAt:			%v\n", asset.GetCreatedAt())
				fmt.Printf("	UpdatedAt:			%v\n", asset.GetUpdatedAt())
				fmt.Printf("	BrowserDownloadURL:	%v\n", asset.GetBrowserDownloadURL())
				fmt.Printf("	State:				%v\n", asset.GetState())
				fmt.Printf("	ContentType:			%v\n", asset.GetContentType())
				fmt.Printf("	DownloadCount:		%v\n", asset.GetDownloadCount())
				fmt.Printf("	NodeID:				%v\n", asset.GetNodeID())
			}
		}

		eblog.Debug(entity.VmBoxEntityName, "Showing all ISO releases. Latest == %s", me.Latest)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *Releases) UpdateReleases() error {

	var rm = make(ReleasesMap)
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.BaseDir == "" {
			p := osdirs.New()
			me.BaseDir = p.AppendToUserConfigDir("iso")
		}

		me.Map = rm

		client := github.NewClient(nil)
		//ctx := context.Background()
		opt := &github.ListOptions{}

		releases, _, err := client.Repositories.ListReleases(context.Background(), "gearboxworks", "gearbox-os", opt)
		if err != nil {
			err = msgs.MakeError(entity.VmBoxEntityName, "can't fetch GitHub releases")
			break
		}

		findFirst := true
		for _, rel := range releases {
			if rel == nil {
				continue
			}

			name := Version(rel.GetName())

			release := Release{
				Version:  name,
				Url:      "",
				Instance: rel,
				channels: me.channels,
			}

			// rm[name].Url/File - Find the first ISO asset.
			for _, asset := range rel.Assets {
				if strings.HasSuffix(asset.GetBrowserDownloadURL(), ".iso") {
					// Return the first ISO found.
					release.Url = asset.GetBrowserDownloadURL()
					release.File = osdirs.AddFilef(me.BaseDir, asset.GetName())
					release.Size = int64(asset.GetSize())
					break
				}
			}

			// rm[name].Version - Copy version name.
			rm[name] = &release

			// rm.Latest - Find first version and select as 'latest'.
			if findFirst {
				me.Latest = &release
				findFirst = false
			}
		}

		//if findFirst == true {
		//	// If we never found a "first", then there was no data.
		//	// So don't update the map.
		//}

		me.Map = rm

		eblog.Debug(entity.VmBoxEntityName, "Fetching ISO releases. Latest == %s", me.Latest)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

var _ hypervisor.Releaser = (*Release)(nil)

type Release struct {
	Version       Version
	File          osdirs.File
	Size          int64
	Url           string
	Instance      *github.RepositoryRelease
	DlIndex       int
	IsDownloading bool

	channels *channels.Channels
}

func (me *Release) GetFilepath() string {
	return me.File
}

const IsoFileNeedsToDownload = 0
const IsoFileIsDownloading = 1
const IsoFileDownloaded = 2

func (me *Release) IsIsoFilePresent() (int, error) {

	var err error
	var ret int
	var stat os.FileInfo

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.File == "" {
			err = msgs.MakeError(entity.VmBoxEntityName, "no Gearbox OS iso file defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			break
		}

		stat, err = os.Stat(me.File)
		if os.IsNotExist(err) {
			err = msgs.MakeError("", "ISO file needs to download from GitHub VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			ret = IsoFileNeedsToDownload
			break
		}

		if me.IsDownloading {
			err = msgs.MakeError("", "ISO file still downloading VmIsoUrl:%s VmIsoFile:%s Percent:%d", me.Url, me.File, me.DlIndex)
			ret = IsoFileIsDownloading
			break
		}

		if stat.Size() != me.Size {
			err = msgs.MakeError("", "ISO file needs to re-download from GitHub VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			ret = IsoFileNeedsToDownload
			break
		}

		ret = IsoFileDownloaded
		me.DlIndex = 100
		eblog.Debug(entity.VmBoxEntityName, "ISO already fetched from '%s' and saved to '%s'", me.Url, me.File)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return ret, err
}

func (me *Release) ShowRelease() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.Instance.Name == nil {
			err = msgs.MakeError(entity.VmBoxEntityName, "no release version specified")
			break
		}

		fmt.Printf("Assets for release:	%v\n", *me.Instance.Name)
		for _, asset := range me.Instance.Assets {
			fmt.Printf("	Name:				%v\n", asset.GetName())
			fmt.Printf("	ID:					%v\n", asset.GetID())
			fmt.Printf("	URL:					%v\n", asset.GetURL())
			fmt.Printf("	Size:				%v\n", asset.GetSize())
			fmt.Printf("	CreatedAt:			%v\n", asset.GetCreatedAt())
			fmt.Printf("	UpdatedAt:			%v\n", asset.GetUpdatedAt())
			fmt.Printf("	BrowserDownloadURL:	%v\n", asset.GetBrowserDownloadURL())
			fmt.Printf("	State:				%v\n", asset.GetState())
			fmt.Printf("	ContentType:			%v\n", asset.GetContentType())
			fmt.Printf("	DownloadCount:		%v\n", asset.GetDownloadCount())
			fmt.Printf("	NodeID:				%v\n", asset.GetNodeID())
		}

		eblog.Debug(entity.VmBoxEntityName, "Showing ISO release for v%s", *me.Instance.Name)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *Release) GetIso() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.File == "" {
			err = msgs.MakeError(entity.VmBoxEntityName, "no Gearbox OS iso file defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			break
		}

		if me.Url == "" {
			err = msgs.MakeError(entity.VmBoxEntityName, "no Gearbox OS iso url defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			break
		}

		var state int
		state, err = me.IsIsoFilePresent()
		if state != IsoFileNeedsToDownload {
			break
		}

		// Start download
		me.DlIndex = 0
		me.IsDownloading = true
		client := grab.NewClient()
		req, _ := grab.NewRequest(me.File, me.Url)
		eblog.Debug("", "downloading ISO from URL %s", req.URL())
		resp := client.Do(req)
		// fmt.Printf("  %v\n", resp.HTTPResponse.Status)
		fmt.Printf("%s VM: Downloading ISO from '%s' to '%s'. Size:%d\n",
			global.Brandname,
			me.Url,
			me.File,
			resp.Size)

		// start UI loop
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-t.C:
				me.DlIndex = int(100 * resp.Progress())
				me.publishDownloadState()
				//fmt.Printf("Downloading '%s' transferred %v / %v bytes (%d%%)\n", me.File, resp.BytesComplete(), resp.Size, me.DlIndex)
				fmt.Printf("%s VM: Downloading ISO - %d%% complete.\r",
					global.Brandname,
					me.DlIndex)

			case <-resp.Done:
				// download is complete
				break Loop
			}
		}

		// check for errors
		if err := resp.Err(); err != nil {
			fmt.Printf("\nDownload failed\n")
			err = msgs.MakeError(entity.VmBoxEntityName, "ISO download failed VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File)
			break
		}
		fmt.Printf("%s VM: Downloaded ISO completed OK.\n",
			global.Brandname,
		)

		eblog.Debug(entity.VmBoxEntityName, "ISO fetched from '%s' and saved to '%s'. Size:%d", me.Url, me.File, resp.Size)
		me.DlIndex = 100
		me.publishDownloadState()
		me.IsDownloading = false
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *Release) publishDownloadState() {

	client := msgs.Address(entity.VmUpdateEntityName)
	state := states.New(client, client, entity.VmBoxEntityName)
	state.SetWant("100%")
	state.SetCurrent(states.State(fmt.Sprintf("%d%%", me.DlIndex)))

	f := msgs.Address(states.ActionUpdate)
	msg := f.MakeMessage(entity.BroadcastEntityName, states.ActionStatus, state.ToMessageText())
	_ = me.channels.Publish(msg)
}

func (me *Release) EnsureNotNil() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("release is nil")
			break
		}
	}

	return err
}
