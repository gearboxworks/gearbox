## Changelog

19c364b (almost?) finalized project api, and work on stack api.
bb5094e .
15dab4e A bit more testing.
a80de3d A little cleanup
0bf22e7 Ability to change service version and ability to remove project stack
c844e65 Accept partial matching for services (e.g. php:7.1 in place of php:7.1.18)
5c92695 Add /.gearbox to .gitignore
cdc1476 Add data/stack-member-options.json
121cd03 Add stack to project mostly working
74b52be Add test/user-home/.gearbox/cache/ to .gitignore
4cc7e65 Added /meta/endpoints and /meta/methods to the API
a927b1f Added BasedirController and Basedirmodel to apimvc and support in config.Config.
b1f0512 Added GET,HEAD,POST for Directories to API.
f28734b Added Gearspec and Authority to the API
232b33a Added MQTT <-> channel proxy
95dc6b7 Added Optional to stack members, and to CacheServer
bd6b185 Added Optional to stack members, and to CacheServer, take 2
ccdd110 Added Run/Stop, Home, Dashboard icons; improved layout for project details
6c7ffaa Added ValuesFunc to API routes definition
e203064 Added Viewing Options panel on Projects master page (not functional yet)
dce997b Added WinConnector
bd77ac0 Added `--no-cache' to global options
188a5da Added `root` and `current` links, changed `item` link to `list` where appropriate
65fc318 Added a @TODO comment for Mick
04e567a Added a GetCacheDir() to host connector.
05a958e Added abilit to add a basedir.  Updated API to support *Status instead of error.
5ab2a7d Added application/json header
abc8a1a Added boilerplate code and build scripts for Admin
5b838e2 Added constants for the stack section of the API
b11c8ad Added error handling and renamed for idiomatic Go.
545be95 Added error handling on failed copy
35ddfa4 Added explicit URL
665380d Added first ValuesFunc in getStackNameValues()
f86d982 Added gearbox.GetFullStackNames()
23bfa41 Added gearbox/assets.go
40f3a36 Added import "text/template" to daemon_windows.go
534ba35 Added in some missing status handling
18b0f68 Added labels and Viewing Options panel to show current settings; Now showing actual basedirs and stacks in filter dropdowns
765f91c Added links for relationships
7389f52 Added logic to use ListCollection.GetRelatedFields()
2704313 Added multiple stacks to options.json, for data during developer/debugging/testing
e4d945e Added notifications at project and stack level;
41c4345 Added objects for Gearbox, Projects, and Config.
ab97266 Added option of using Lorca or Webview for admin viewer.
2c9857b Added options
102218c Added options.
41cf45a Added proper tabular view for projects; some restructuring for reuse
3ead679 Added response message for add basedir
9181322 Added some printfs
1bd9348 Added support for http://127.0.0.1:9999/basedirs/new to add a new basedir. Likely not robust yet.
158d718 Added terminal W/H.
41b9d21 Added test.HttpClient() needed to test API
60452a5 Added tests for caching
a413341 Added the latest Mac executable
0816be8 Added type aliases for status.Status and apimodeler.* types.
0a01520 Added viewing filter to limit the list of projects by program
bb4b431 Adding in Webview and a `gearbox admin` command.
59f72e5 Adding stashes as patches
8c780ca Admin - work in progress
e89ac22 All uncommented tests not passing again
9a5a88d Allow deselecting services that are optional; Sanitize project path;
15a03a4 And .DS_Store
0e6bbe7 Attempt to address CORS with middlewear
82d1c9d Attempt to fix Windows compile issue, take 2.
fa76548 Attempt to fix Windows compile issue, take 3.
a0e61e4 Attempt to fix Windows compile issue, take 4.
14f0073 Attempt to fix Windows compile issue.
0d17b1d Basedir Add and Update API fixed to work with JSON:API spec: https://jsonapi.org/format/1.1/#crud
1210aa4 Better PTY handling.
6b22251 Bug fixes
13d26f0 CORS for API
b2583f8 Change "label" => "name", "short_label" => "label"
8a4de0d Change from using hostname to using path for project dir.
a4b6393 Change to use admin_dist.go instead of assets.go
595aebc Changed "role_services" to "services" in online `gears.json`
97a1cb7 Changed ESLint config to Standard
ed6beff Changed `Status.Success` to `.Failed` so `IsSuccess()` is `true` for an empty `Status{}`
469b635 Changed `gearbox.json` to use `"schema"` instead of `"json"`, and dropped box `"version"` (might add back later.)
b3affca Changed appdist to app/dist
431bc46 Changed errors responses to be consistent with other responses.
e61a906 Changed getStackName() to recognize and return authority prefixed to stackname if found
6febb13 Changed model types to plural because @reststate/vuex is not flexible enough to support singular types.
0f8094f Changed to returning Status as non-pointers to ensure they will always be available (e.g. non-nil)
788e8c4 Changed to using ResourceVarName instead of string for RequestContext URL segment params.
41946b2 Changes after working with Saru about first admin milestone.
47bd8c1 Checkin #1 before merge.
5fca6cf Cleaning Up #2
18538c3 Cleaning up #1
4a46f5a Cleaning up Response constants
38f4d71 Cleanup update #3
d09907c Clear up HelpfulError edge cases that caused bugs
d09710f Collapse project [location]/[add-note] controls when ESC key is hit; some styling tweaks;
635e57f Colours!
9b58d7c Commit a gear change: "org" => "orgname"
79da659 Commit a gear change: remove accidental inclusion of "program" and duplicates of "default" for "services"
e505fff Commit after merge
8bb6312 Commit after merge
b5d8411 Completed UI for adding/removing project stacks
eca44a4 Completed basic testing for apimodeler
4a838e5 Completed basic testing for apimodeler
5287d5e Delete exports
934e5c4 Deleted cardinality. Seemed like a good idea at the time.
d59c4c2 Deleted data/stack-member-options.json and added options.json
62c9ae9 Do not allow switching services while the projects is in a running state; Show notification when selected service version is different from spec
3b00136 Do not pass nickname and id when calling basedir create
311daee Doh!  Wrong format for Service versions.
c5bf57c Dropped vm from Api Service name
42554d6 Enable ISO download from GitHub
aefd6b0 Enabled project Stop/Run functionality
4b1169f EventBroker
12727a9 Externalize cache package.
37ae296 Extract ProjectCfg to a separate repo.
0925796 Extracted OsSupport as go-osbridge in a separate repo.
0d57539 File delete method
e654161 File path check.
b18fe50 File perms
3a5663b File perms.
a1a385b Final Heartbeat integration.
9cc0889 Final commit before merge.
f4f57af Finally got projects and project details API routes working again
a395198 Finetuned /stacks, added list, item and common link types.
6259937 Finetuning of stacks in API response for projects
e3116c6 First pass (untested) of setting related items
50650c3 Fix API for add and update to drop /new and use PATCH vs. PUT.
8937f03 Fix POST add for Directories to API.
dc88c33 Fix ability to create ~/.gearbox and to write config.json
63d322a Fix bug that entered into AddBaseDir
ecaa5d2 Fix build tags
da12def Fix compilation directives in os_support files
65377a7 Fix error that made /projects empty
23566df Fix for move from `appdist` to `app/dist`
1c298fa Fix issue stoping loading of projects.
2414f30 Fix mistaken search-and-replace of Gearbox => Parent
ed13f62 Fix too many Vary headers in API
f8fd06a Fix unfortunate rename of 'Gearbox' to 'Parent' in text.
60725f3 Fixed Windows-specific compilation error in osbridge_windows.go
78159a5 Fixed bug in api.GetRequestContext() that did not trim slashes on the URL path
3d320dc Fixed catch-22 of config file missing or can't load.
deb104d Fixed catch-22 of config file missing or can't load.
06892ab Fixed logic that did not append '.local' to project path to get a valid project hostname.
9fb0dd0 Fixed one copy-paste and one lack of error handling bugs.
559bdef Fixed short_label response.
964ff57 Fixed the removal of stack members from stack API response.
22755ad Fixes to basedir API.
7f9b8a3 Fixes to return full Resource Object on add via API.
8dfd509 Fixes to return full Resource Object on update via API.
d699a39 Fixing Project Candidates to have the correct response
27868ea Fleshed out MacOs/Win/Linux connectors.
7492f6d For convenience, added assets in dist/js
43db848 Further refactoring of HostApi
ebcafe6 Got /stacks/:authority/:stackname serving
2c51873 Got project resource working in API.
7e3455d Got rid of all other 'only' packages except github.com/gearboxworks/go-status/only
04805be Gotta .gitignore .idea cause it causes too much trouble when there is a merge conflict.
0030f24 Highlight non-default viewing options (i.e. active filters); Do not allow unsellecting all project states
b6d2824 Implement Gearbox.IsDebug()
3146c91 Implement a simple JSON cache-to-file package.
a50357a Implemented "projects-with-details" API resource.
61026dc Implemented 'options.json' downloaded, v1
b165235 Implemented /services (finally.)
2d8c973 Implemented API delete generically and for the basedir use-case. Also renamed apimodeler package to apiworks.
99b1ede Implemented Update and Delete for BaseDirs
bc31ce6 Implemented architecture for global options and an IsDebug global option.
8c43817 Implemented edit form for Project Details, POSTing project data to the API
af37929 Implemented error logging and cleaned up CLI output
979b27e Implemented project stacks as tags that expand to mini-cards when clicked
370133d Implemented sidebar navigation, now Project and Stack fields are rendered
f89a6c8 Implemented stack, stack members, stack member options in the API.
a779807 Improve the way the closest match for a service version is resolved and displayed
4d151fa Improved UX for expanding/collapsing project details section
13cbde3 Improved card layout on Projects master page
1ed54f9 Improved project candidate output.
b19e6ff Improved state management on Preferences page
a9a63b7 Improved styling and methods for adding/changing/deleting basedirs
23fc876 Initial
6a06c17 Initial
abad344 Initial code layout.
b6f1fc1 Initial commit
81d30a4 Initial commit (includes GoLand configuration.)
23257ce Introduced stack action icons; now hostnames can start with a digit
aff8721 Major refactor
ba88222 Major refactor #2
b7036b0 Many many changes.
c3d0fe6 Mega update prioer to merge.
d5e5eea Mega update.
cbd9b95 Merge branch '1-lowlevel-vm-control'
dee22b1 Merge branch '2-VM-Config-Creation'
02a9879 Merge branch 'jsonapi'
6b49e48 Merge branch 'master' into 7-Host-Hearbeat-Part2
66f5745 Merge branch 'master' of github.com:gearboxworks/gearbox
e1287be Merge branch 'master' of github.com:gearboxworks/gearbox
4551ec7 Merge branch 'master' of github.com:gearboxworks/gearbox
8bbfc2c Merge branch 'master' of https://github.com/gearboxworks/gearbox
a5f6cca Merge branch 'master' of https://github.com/gearboxworks/gearbox
4440c72 Merge branch 'master' of https://github.com/gearboxworks/gearbox
14c917f Merge branch 'master' of https://github.com/gearboxworks/gearbox
12cbbf4 Merge branch 'master' of https://github.com/gearboxworks/gearbox
9a678b2 Merge branch 'master' of https://github.com/gearboxworks/gearbox
6f5d6b5 Merge branch 'master' of https://github.com/gearboxworks/gearbox
e8bb753 Merge branch 'master' of https://github.com/gearboxworks/gearbox
2856c8f Merge in old heartbeat.
939f5cb Merge pull request #13 from gearboxworks/Merge-7-Host-Hearbeat
c1c9199 Merge pull request #20 from gearboxworks/7-Host-Hearbeat-Part2
34c06de Merge remote-tracking branch 'origin/master'
2b55e17 Merge remote-tracking branch 'origin/master'
d4db6b6 Merge remote-tracking branch 'origin/master'
2c70921 Merge remote-tracking branch 'origin/master'
0e7cb5a Merge remote-tracking branch 'origin/master'
8e73f1f Merge remote-tracking branch 'origin/master'
239ee60 Merge remote-tracking branch 'origin/master'
3a99094 Merge remote-tracking branch 'origin/master'
f4f2b9f Merge remote-tracking branch 'origin/master'
eb1e2ee Merge remote-tracking branch 'origin/master'
864fef7 Merged [Edit Note] and a list of note icons into a single component; now allowing only a single note per project
b71a3bc Merging all command line option into one command line app.
c1af688 Milestone #2: VM GoLang status confirmation.
c5884f4 Milestone #3: SSH
21bf633 Milestone #4: VM and SSH functional
34df7ce Milestone update
4ac9ef7 Milestone update #2
021fd8d Milestone: correct VM start/stop
877fe25 Minor changes to support the current gears.json schema
cfe9665 Minor changes.
d5683b9 Minor changes.
ebb6e9b Minor fix to allow changing service version in project stack
1af878c Mock objects for Gearbox and HostConnector and initial testing for API responses.
b540c36 Modified Admin build script to copy admin/gears.json to /dist
847b7b7 Modified Admin to work with the new API responses
7beacc7 Modified build script to copy all files from /assets to /admin/dist
caf6b92 More implementation of JSON:API `relationships` and `included`
fe603c3 More of /stacks working again.
f1c1d4b More potential fixes for Windows not loading projects
1a2d5f7 More work on the admin console
4ae6cd5 More work-in-progress on refactoring to support JSON:API
137f6b1 Move refactoring
15eec9c Move towards standardizing interface for HostApi handlers
0bc7a71 Move towards standardizing interface for HostApi handlers
297d5a9 Moved Basedir.Nickname to UUID.
c44bcd1 Moved options.json and created test app to debug it.
b3bff63 Moved package up a level.
d319833 Moved stack selector inside project details section
8e2c94d Mutable mutex maps.
94e4a45 NFS fixups.
ece4ba5 NIC changes.
fa7c3f3 Nada
c79096c Nascent tabular view for Projects master page
8b444a5 New messaging system
f8d9097 Now a warning message is shown in Admin when server does not respond (and retrying 5 times before giving up)
333b276 Now actually filtering projects by the specified criteria
95ba55c Now collapsible stack sections are rendered based on project settings; Version numbers are displayed in the header of collapsible section
0a4047c Now rendering accordion with stack services
ff54184 Now stack action icons are removed with a transition when stopping a project; added Folder icon to project details section
8c9014b OVA asset
739710c On cards view, show project stack mini-card readily expanded when there is only one stack on a project
2b48eb8 Ongoing refactoring, and implemented generic API and basedir API support for PUT to update.
71c6f84 Output the actual error when RestoreAsset() fails.
1c0e90e Partial implementation of JSON:API `relationships` and `included`
06303f1 Passing nickname with basedir update, generating id for basedir create
9c06f47 Path checking.
0433fb1 Post merge fixups.
37ed61a Post merge fixups.
0a18582 Potential fix for Windows not loading projects
c8b044d Preselect proper service version in a droplist, fix styling
ee93508 Proof-of-concept for Api
aa75315 Proof-of-concept for Api
e32cd29 Proof-of-concept for plugins
948b602 Proof-of-concept working for passing data from Go to JS.
bc7e1c9 Properly added new /admin files
dbd1196 Properly copy files from /assets to /admin/dist
0e6adb0 Pushing non-dev assets. But they will be regenerated.
bdc19d1 Reduce shutdown time.
1445c0b Refactored HostAPI into multiple files.
edd20b4 Refactored api.WireRoutes() to break out the different routes, and renamed 'ja' package to 'jsonapi'
9286a87 Refactored from `options.json` to `gears.json`, etc. Refactored more error handling. Redirected trailing slash in endpoints.
9f7d932 Refactored preferencess page; Added [Copy to clipboard] and [Open in file manager] buttons for basedirs; Added updating spinners; Still no checking if dirs actually exist
9a67446 Refactored to eliminate HelpfulError and replace with stat.Status everywhere.
8bd3983 Refactored to get /stacks partially working again.
b35fb33 Refactored to simplify implement additional Api routes
387409a Refactored to use our own apimodeler.Context struct that embeds echo.Context instead of apimodeler.Contexter interface.
899bf6d Refactoring and Project Services + Project Aliases endpoint resource.
cb737df Refactoring and Project Services + Project Aliases endpoint resource.
fe06763 Refactoring to limit Response methods to HostApi, and to implement ProjectStacksResponse
bd091c0 Refactoring to reduce coupling between jsonapi and apimodeler+api
d9e136a Regular commit.
c47fef5 Remove 'bin/gearbox-api' from repo
bd817b7 Remove dangling ref to github.com/wplib/go-lib
0620ffd Remove go.mod from the root.
869b124 Remove guff
959a5c1 Remove old spinner.
24ed224 Remove omit_empty for ProjectStackItem.ServiceId
6c07ffa Remove stack members from stack API response.
9d79fa1 Removed only package from the embedded projectcfg.
6a71213 Removed unused code.
88296ce Removed.
44ebaa9 Rename /gearbox/gearbox to /gearbox/app
aefc70b Renamed "name" to "program"
11d86db Renamed /gearbox/gearbox to /gearbox/app
0a16bb8 Renamed API links from  "project" and "stack" to "project-details" and "stack-details", respectively.
a544cff Renamed Api to HostApi, moved into the gearbox.Admin() method.
1837fb9 Renamed GearboxArgs to just Args
69c5e91 Renamed NewFailedStatus() to NewFailedStatus(), added NewErrorStatus()
6d3ac8e Renamed all things from "vm" to "box"
e478d4b Renamed apibuilder to modeler, added more tests. Fixed some typos in repo.
f30baa4 Renamed links. Added stub for "basedir-delete" api.
441b086 Renamed options.json to gears.json
f404961 Replace getlantern/systray with gearboxworks/go-systray
9f7871d Replaced Element UI with Bootstrap Vue; replaced sidebar with cards
bef2e63 Replaced Vuetify component lib with ElementUI, now MSHTML seems to work; Now fetching and rendering names of actual projects/stacks from the API. Attempted to add CORS rules to API response headers to be able to test Admin from Chrome but unsuccessfully (left some commented-out code in api.go).
82b3d2b Replaced static info on project cards with editable forms;
0aff36b Representing services with icons (to brighten up project stack section)
b96bfe6 Restore `Gearbox Sites` as primary base dir on Windows
3cb4727 Restored ability to add stacks to projects
5c723b2 Restored support for optional services
5b20918 Restructure Basedir struct dropping BoxDir property and renaming HostDir property to Basedir. Also renamed 'primary' basedir to 'default'
3adabfc Restructured options schema
54269a0 Restructured web server so it does not leave orphans.
cc93c8c Reverted the hiding of `status.success` property.
20b9496 Reverts basedir IDs using UUIDs and instead sets them based on the last segment of the basedir while ensuring they are unique.
f4b41ee Reworked project cards; added sensible submit buttons and update spinners for every modifiable field; updated tabular view with the new components;
59d0c2d SSH config
4955818 Save VM configs as JSON.
ceba57e Shortened labels for Project State filters; removed directory modification popup; minor css tweaks
8c8c591 Show NO location filter and no badge when there is only one basedir; Fixed as small toggling issue in project state filter
dc2cda4 Show project [location]/[add-stack]/[add-note] controls inline, autofocus when expanded
f861673 Show project details only after project title is clicked
b1db051 Show spinners while switching project state and while loading project details
4bf8eda Small changes
34c7b11 Some tweaks to make Viewing Options look better (less broken) in smaller window
fed030f Stupid Git keeps stuffing things up completely.
6ec4016 TTY changes.
90702e0 Task process management
6c8ceda TestGlobalOptions passing again.
b6327bb Update #2
bbbee47 Update #3.
587773b Update #4.
ca93fb7 Update #5.
456a107 Update #6.
3fa9071 Update .goreleaser.yml
6c98953 Update card caption when a different service program/version is selected
c8b4359 Update dist.go
88ac9e3 Update dist.go
a48ae13 Update gears.json, in three places :-(
753a08b Update go.sum
2af48fd Update heartbeat.go
17eacdc Update menus.go
bbc2cb0 Update ssh.go
5086eab Update the .paw file
77d770c Update virtualbox.go
43a003a Update virtualbox.go
4890a6b Update virtualbox.go
cbc2e4e Update vm.go
8cd5ba3 Update vm.go
f2775cf Updated Api package to return a Resource Identifier Object for an add. Now returns proper status code on an add. Validated when basedir "id" posts does not equal "nickname" or set nickname from id.
542ab12 Updated JS packages
c2b6449 Updated Lorca to v0.1.7 (because v0.1.6 was failing after upgrading to go v1.12.6)
c225f00 Updated most of the templates and rendering functions to use the new API (interactive features still require more work)
6c85722 Updated to newer version of go-status, again
cf7ed80 Updated to newer version of go-status, redux
bae2953 Updated to newer version of go-status.
49e005c Updated.
0894305 Updated.
2dee174 Updates go-status and fixed resultant breakage
d490109 Updates to support Basedir validation, changes to README
77a7fb4 Use basedir data from the API
3c5e932 Use filepath.FromSlash to have portable directory separator
689bffc Use gearbox.local config.
4deb660 Use gearbox.local host
da0f972 Use some scaling animation to indicate service change
2031f8f Use stack mini-cards in cards view mode of the projects list
88e828a Use temp var for httpStatus
a2011c1 Work in progress
0bfc91d Work in progress for the API's Project resource
17ffb28 Work-in-progress for adding a new stack to a project
d76d716 Work-in-progress on HTTP testing the responses from the API
b19233d Work-in-progress on refactoring to support JSON:API
d825bc2 cleanup
9054438 const name changes.
8e0e181 fn name changes.
0829f23 hostonlyif changes
warning: refname '0.5.1' is ambiguous.
57ad149 removed commented lines from api.go
ab56f92 zero net fix.
