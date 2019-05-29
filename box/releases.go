package box

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	"github.com/cavaliercoder/grab"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"os"
	"strings"
	"time"
)


type Release github.RepositoryRelease
type ReleasesMap map[string]github.RepositoryRelease

type Releases struct {
	Map ReleasesMap
	Latest	string
}

type ReleaseSelector struct {
	// These are considered to be AND-ed together.
	FromDate        time.Time
	UntilDate       time.Time
	SpecificVersion string
	RegexpVersion   string
	Latest			*bool
}


func (me *Releases) ShowReleases() (status.Status) {
	var sts status.Status

	for range only.Once {
		sts = EnsureReleasesNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Latest release: %v\n\n", me.Latest)

		for _, release := range me.Map {
			fmt.Printf("Assets for release:	%v\n", release.GetName())
			fmt.Printf("UploadURL: 			%v\n", release.GetUploadURL())
			fmt.Printf("ZipballURL: 			%v\n", release.GetZipballURL())
			fmt.Printf("TarballURL: 			%v\n", release.GetTarballURL())
			fmt.Printf("Body: 				%v\n", release.GetBody())
			fmt.Printf("AssetsURL: 			%v\n", release.GetAssetsURL())
			fmt.Printf("URL: 				%v\n", release.GetURL())
			fmt.Printf("HTMLURL:				%v\n", release.GetHTMLURL())

			for _, asset := range release.Assets {
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

		sts = status.Success("%s VM - Showing all ISO releases. Latest == %s.\n", global.Brandname, me.Latest)
	}

	return sts
}


func (release *Release) ShowRelease() (status.Status) {
	var sts status.Status

	for range only.Once {
		sts = EnsureReleaseNotNil(release)
		if is.Error(sts) {
			break
		}
		if release.Name == nil {
			sts = status.Fail().
				SetMessage("no release version specified").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		fmt.Printf("Assets for release:	%v\n", *release.Name)
		for _, asset := range release.Assets {
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

		sts = status.Success("%s VM - Showing ISO release for v%s.\n", global.Brandname, *release.Name)
	}

	return sts
}


func (me *Box) GetReleases() (Releases, status.Status) {

	var r Releases
	var rm = make(ReleasesMap)
	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		client := github.NewClient(nil)
		ctx := context.Background()
		opt := &github.ListOptions{}
		releases, _, err := client.Repositories.ListReleases(ctx, "gearboxworks", "gearbox-os", opt)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("can't fetch GitHub releases").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		findFirst := true
		for _, release := range releases {
			if findFirst {
				r.Latest = release.GetName()
				findFirst = false
			}

			// Convert to a map.
			rm[release.GetName()] = *release
		}

		r.Map = rm

		sts = EnsureReleasesNotNil(&r)
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s VM - Fetching ISO releases. Latest == %s.\n", global.Brandname, r.Latest)
	}

/*
type Release struct {
	Name string
	UploadURL string
	ZipballURL string
	TarballURL string
	Body string
	AssetsURL string
	URL string
	HTMLURL string
	Name string
    Assets
}
type Releases []Release

type Asset struct {
      Name
      ID
      URL
      Size
      CreatedAt
      UpdatedAt
      BrowserDownloadURL
      State
      ContentType
      DownloadCount
      NodeID
}
type Assets []Asset


*/


/*
   Assets for release:	0.5.0
   UploadURL: 			https://uploads.github.com/repos/gearboxworks/gearbox-os/releases/17531887/assets{?name,label}
   ZipballURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/zipball/0.5.0
   TarballURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/tarball/0.5.0
   Body: 				0.5.0 pre-release
   AssetsURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/releases/17531887/assets
   URL: 				https://api.github.com/repos/gearboxworks/gearbox-os/releases/17531887
   HTMLURL:				https://github.com/gearboxworks/gearbox-os/releases/tag/0.5.0
   foo: 				0.5.0
   Name:				gearbox-0.5.0.iso
   ID:					12825393
   URL:					https://api.github.com/repos/gearboxworks/gearbox-os/releases/assets/12825393
   Size:				67108864
   CreatedAt:			2019-05-23 02:37:48 +0000 UTC
   UpdatedAt:			2019-05-23 02:42:56 +0000 UTC
   BrowserDownloadURL:	https://github.com/gearboxworks/gearbox-os/releases/download/0.5.0/gearbox-0.5.0.iso
   State:				uploaded
   ContentType:			application/octet-stream
   DownloadCount:		0
   NodeID:				MDEyOlJlbGVhc2VBc3NldDEyODI1Mzkz
 */

/*
			fmt.Printf(`
				release.ID=%v
				release.TagName=%v
				release.TargetCommitish=%v
				release.Name=%v
				release.Body=%v
				release.Draft=%v
				release.Prerelease=%v
				release.CreatedAt=%v
				release.PublishedAt=%v
				release.URL=%v
				release.HTMLURL=%v
				release.AssetsURL=%v
				release.Assets=%v
				release.UploadURL=%v
				release.ZipballURL=%v
				release.TarballURL=%v
				release.Author=%v
				release.NodeID=%v\n`,
				release.ID,
				release.TagName,
				release.TargetCommitish,
				release.Name,
				release.Body,
				release.Draft,
				release.Prerelease,
				release.CreatedAt,
				release.PublishedAt,
				release.URL,
				release.HTMLURL,
				release.AssetsURL,
				release.Assets,
				release.UploadURL,
				release.ZipballURL,
				release.TarballURL,
				release.Author,
				release.NodeID,
				)
*/

	return r, sts
}


/*
Updates the following:
   me.VmIsoVersion    string
   me.VmIsoFile       string
   me.VmIsoUrl 		string
   me.VmIsoRelease    Release
*/
func (me *Box) SelectRelease(selector ReleaseSelector) (status.Status) {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		releases, sts := me.GetReleases()
		if is.Error(sts) {
			break
		}

		me.VmIsoInfo = Release(releases.Map[releases.Latest])
		me.VmIsoVersion = *me.VmIsoInfo.Name
		for _, asset := range me.VmIsoInfo.Assets {
			if strings.HasSuffix(asset.GetBrowserDownloadURL(), ".iso") {
				// Return the first ISO found.
				if (me.VmIsoUrl == "") && (me.VmIsoFile == "") {
					me.VmIsoUrl = asset.GetBrowserDownloadURL()
					me.VmIsoFile = me.VmIsoDir + "/" + asset.GetName()
				}

				break
			}
		}

		sts = EnsureReleaseNotNil(&me.VmIsoInfo)
		if is.Error(sts) {
			break
		}

		sts = status.Success("%s VM - Selected ISO release %s.\n", global.Brandname, me.VmIsoVersion)
	}

	return sts
}

/*
type ReleaseAsset struct {
	ID                 *int64     `json:"id,omitempty"`
	URL                *string    `json:"url,omitempty"`
	Name               *string    `json:"name,omitempty"`
	Label              *string    `json:"label,omitempty"`
	State              *string    `json:"state,omitempty"`
	ContentType        *string    `json:"content_type,omitempty"`
	Size               *int       `json:"size,omitempty"`
	DownloadCount      *int       `json:"download_count,omitempty"`
	CreatedAt          *Timestamp `json:"created_at,omitempty"`
	UpdatedAt          *Timestamp `json:"updated_at,omitempty"`
	BrowserDownloadURL *string    `json:"browser_download_url,omitempty"`
	Uploader           *User      `json:"uploader,omitempty"`
	NodeID             *string    `json:"node_id,omitempty"`
}

type RepositoryRelease struct {
	ID              *int64         `json:"id,omitempty"`
	TagName         *string        `json:"tag_name,omitempty"`
	TargetCommitish *string        `json:"target_commitish,omitempty"`
	Name            *string        `json:"name,omitempty"`
	Body            *string        `json:"body,omitempty"`
	Draft           *bool          `json:"draft,omitempty"`
	Prerelease      *bool          `json:"prerelease,omitempty"`
	CreatedAt       *Timestamp     `json:"created_at,omitempty"`
	PublishedAt     *Timestamp     `json:"published_at,omitempty"`
	URL             *string        `json:"url,omitempty"`
	HTMLURL         *string        `json:"html_url,omitempty"`
	AssetsURL       *string        `json:"assets_url,omitempty"`
	Assets          []ReleaseAsset `json:"assets,omitempty"`
	UploadURL       *string        `json:"upload_url,omitempty"`
	ZipballURL      *string        `json:"zipball_url,omitempty"`
	TarballURL      *string        `json:"tarball_url,omitempty"`
	Author          *User          `json:"author,omitempty"`
	NodeID          *string        `json:"node_id,omitempty"`
}

 */

/*
Data returned:

release.ID=0xc000289538
release.TagName=0xc0002964c0
release.TargetCommitish=0xc0002964d0
release.Name=0xc0002964e0
release.Body=0xc000296770
release.Draft=0xc00028955b
release.Prerelease=0xc00028957d
release.CreatedAt=2019-05-23 02:34:10 +0000 UTC
release.PublishedAt=2019-05-23 02:43:04 +0000 UTC
release.URL=0xc000296470
release.HTMLURL=0xc0002964a0
release.AssetsURL=0xc000296480
release.Assets=[github.ReleaseAsset{
	ID:12825393,
	URL:"https://api.github.com/repos/gearboxworks/gearbox-os/releases/assets/12825393",
	Name:"gearbox-0.5.0.iso",
	State:"uploaded",
	ContentType:"application/octet-stream",
	Size:67108864,
	DownloadCount:0,
	CreatedAt:github.Timestamp{2019-05-23 02:37:48 +0000 UTC},
	UpdatedAt:github.Timestamp{2019-05-23 02:42:56 +0000 UTC},
	BrowserDownloadURL:"https://github.com/gearboxworks/gearbox-os/releases/download/0.5.0/gearbox-0.5.0.iso",
	Uploader:github.User{
		Login:"MickMake",
		ID:17118367,
		NodeID:"MDQ6VXNlcjE3MTE4MzY3",
		AvatarURL:"https://avatars0.githubusercontent.com/u/17118367?v=4",
		HTMLURL:"https://github.com/MickMake",
		GravatarID:"",
		Type:"User",
		SiteAdmin:false,
		URL:"https://api.github.com/users/MickMake",
		EventsURL:"https://api.github.com/users/MickMake/events{/privacy}",
		FollowingURL:"https://api.github.com/users/MickMake/following{/other_user}",
		FollowersURL:"https://api.github.com/users/MickMake/followers",
		GistsURL:"https://api.github.com/users/MickMake/gists{/gist_id}",
		OrganizationsURL:"https://api.github.com/users/MickMake/orgs",
		ReceivedEventsURL:"https://api.github.com/users/MickMake/received_events",
		ReposURL:"https://api.github.com/users/MickMake/repos",
		StarredURL:"https://api.github.com/users/MickMake/starred{/owner}{/repo}",
		SubscriptionsURL:"https://api.github.com/users/MickMake/subscriptions"
		},
	NodeID:"MDEyOlJlbGVhc2VBc3NldDEyODI1Mzkz"
	}]
release.UploadURL=0xc000296490
release.ZipballURL=0xc000296760
release.TarballURL=0xc000296750
release.Author=github.User{Login:"MickMake", ID:17118367, NodeID:"MDQ6VXNlcjE3MTE4MzY3", AvatarURL:"https://avatars0.githubusercontent.com/u/17118367?v=4", HTMLURL:"https://github.com/MickMake", GravatarID:"", Type:"User", SiteAdmin:false, URL:"https://api.github.com/users/MickMake", EventsURL:"https://api.github.com/users/MickMake/events{/privacy}", FollowingURL:"https://api.github.com/users/MickMake/following{/other_user}", FollowersURL:"https://api.github.com/users/MickMake/followers", GistsURL:"https://api.github.com/users/MickMake/gists{/gist_id}", OrganizationsURL:"https://api.github.com/users/MickMake/orgs", ReceivedEventsURL:"https://api.github.com/users/MickMake/received_events", ReposURL:"https://api.github.com/users/MickMake/repos", StarredURL:"https://api.github.com/users/MickMake/starred{/owner}{/repo}", SubscriptionsURL:"https://api.github.com/users/MickMake/subscriptions"}
release.NodeID=0xc0002964b0
 */


func (me *Box) GetIso() (status.Status) {

	// me.VmIsoFile
	// me.VmIsoReleaseUrl
	// me.VmIsoReleases

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		if me.VmIsoFile == "" {
			sts = status.Fail().
				SetMessage("no Gearbox OS iso file found").
				SetAdditional("VmIsoUrl:%s VmIsoFile:%s", me.VmIsoUrl, me.VmIsoFile).
				SetData(me.VmIsoFile).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		sts = me.IsIsoFilePresent()
		if is.Success(sts) {
			break
		}

		if me.VmIsoUrl == "" {
			sts = status.Fail().
				SetMessage("no Gearbox OS iso url found").
				SetAdditional("VmIsoUrl:%s VmIsoFile:%s", me.VmIsoUrl, me.VmIsoFile).
				SetData(me.VmIsoFile).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		client := grab.NewClient()
		req, _ := grab.NewRequest(me.VmIsoFile, me.VmIsoUrl)

		// Start download
		fmt.Printf("Downloading %v...\n", req.URL())
		resp := client.Do(req)
		fmt.Printf("  %v\n", resp.HTTPResponse.Status)
		fmt.Printf("%s VM - ISO fetching from '%s' and saved to '%s'. Size:%s.\n", global.Brandname, me.VmIsoUrl, me.VmIsoFile, resp.Size)

		// start UI loop
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

		Loop:
			for {
				select {
				case <-t.C:
					me.VmIsoDlIndex = int(100*resp.Progress())
					fmt.Printf("File '%s' transferred %v / %v bytes (%d%%)\n", me.VmIsoFile, resp.BytesComplete(), resp.Size, me.VmIsoDlIndex)

				case <-resp.Done:
					// download is complete
					break Loop
				}
			}

		// check for errors
		if err := resp.Err(); err != nil {
			sts = status.Wrap(err).
				SetMessage("iso download failed").
				SetAdditional("VmIsoUrl:%s VmIsoFile:%s", me.VmIsoUrl, me.VmIsoFile).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		fmt.Printf("Download saved to ./%v \n", resp.Filename)

		sts = status.Success("%s VM - ISO fetched from '%s' and saved to '%s'. Size:%d.\n", global.Brandname, me.VmIsoUrl, me.VmIsoFile, resp.Size)
		me.VmIsoDlIndex = 100
	}

	return sts
}

func (me *Box) IsIsoFilePresent() status.Status {

	var sts status.Status

	_, err := os.Stat(me.VmIsoFile)
	if err == nil {
		me.VmIsoDlIndex = 100
		sts = status.Success("%s VM - ISO already fetched from '%s' and saved to '%s'.\n", global.Brandname, me.VmIsoUrl, me.VmIsoFile)

	} else {
		sts = status.Wrap(err).
			SetMessage("can't download iso release from GitHub").
			SetAdditional("VmIsoUrl:%s VmIsoFile:%s", me.VmIsoUrl, me.VmIsoFile).
			SetData("").
			SetCause(err).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("%s VM - ISO not downloaded from '%s' and saved to '%s'.\n", global.Brandname, me.VmIsoUrl, me.VmIsoFile),
			Help:    help.ContactSupportHelp(), // @TODO need better support here
		})
	}

	return sts
}

func EnsureReleaseNotNil(rm *Release) (sts status.Status) {
	if rm == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    VmStateUnknown,
		})
	}

	return sts
}


func EnsureReleasesNotNil(rm *Releases) (sts status.Status) {
	if rm == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    VmStateUnknown,
		})
	}

	return sts
}

