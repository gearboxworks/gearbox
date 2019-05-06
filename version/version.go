package version

import (
	"fmt"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"regexp"
	"strconv"
	"strings"
)

const (
	verMajor = iota
	verMinor
	verPatch
	verPrerelease
	verMetadata
	verRelease
)

type MajorMinor string
type Metadata string
type Revision string
type Prerelease string
type Digits string

type Versioner interface {
	Parse(ver string) (sts status.Status)
	String() string
	GetRaw() string
	GetMajor() string
	GetMinor() string
	GetPatch() string
	GetPrerelease() string
	GetMetadata() string
	GetMajorMinor() string
	GetVersion() string
	GetRevision() string
}

type Version struct {
	raw        types.Version
	Major      Digits        `json:"major,omitempty"`
	Minor      Digits        `json:"minor,omitempty"`
	Patch      Digits        `json:"patch,omitempty"`
	Prerelease Prerelease    `json:"prerelease,omitempty"`
	Metadata   Metadata      `json:"metadata,omitempty"`
	Revision   Revision      `json:"revision,omitempty"`
	Status     status.Status `json:"-"`
}

func NewVersion() *Version {
	return &Version{}
}

func (me *Version) ParseString(gearid string) (sts status.Status) {
	return me.Parse(types.Version(gearid))
}
func (me *Version) Parse(ver types.Version) (sts status.Status) {
	var msg, hlp string
	newBuf := true
	parts := strings.Split(string(ver), "~")
	tmp := Version{raw: ver}
	re := regexp.MustCompile("[^A-Za-z0-9.-]")
	for range only.Once {
		if len(parts) >= 2 {
			if len(parts[1]) == 0 {
				msg = fmt.Sprintf("revision following '~' in '%s' is empty", ver)
				break
			}
			if parts[1][0] != 'r' {
				msg = fmt.Sprintf("revision following '~' in '%s' must begin with 'r'", ver)
				break
			}
			tmp.Revision = Revision(parts[1][1:])
			_, err := strconv.Atoi(string(tmp.Revision))
			if err != nil {
				msg = fmt.Sprintf("revision following '~r' in '%s' must be and integer", ver)
				break
			}
		}
		var s byte
		var i int
		pos := -1
		buf := make([]byte, 0)
		var done bool
		for i, s = range []byte(parts[0]) {
			if newBuf {
				newBuf = false
				pos, msg = tmp.captureMMP(ver, pos, buf)
				if status.IsError(sts) {
					break
				}
				buf = buf[0:0]
			}
			switch s {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				buf = append(buf, s)
			case '.':
				newBuf = true
			case '-', '+':
				done = true
				break
			default:
				msg = fmt.Sprintf("non-integer %s version in '%s'",
					[]string{"major", "minor", "patch"}[pos],
					ver,
				)
				break
			}
			if msg != "" || done {
				break
			}
		}
		if msg != "" {
			break
		}
		_, msg = tmp.captureMMP(ver, pos, buf)
		if msg != "" {
			break
		}
		metadata := Metadata("")
		sharedHelp := "can only contain letters, numbers, periods or dashes. See the SemVer docs for more: https://semver.org/"
		if s == '-' {
			prerelease := ver[i+1:]
			idx := strings.Index(string(prerelease), "+")
			if idx != -1 {
				metadata = Metadata(prerelease[idx+1:])
				prerelease = prerelease[0:idx]
				s = '+'
			}
			if re.MatchString(string(prerelease)) {
				msg = fmt.Sprintf("pre-release in '%s' is invalid semver", ver)
				hlp = fmt.Sprintf("pre-release %s#spec-item-10", sharedHelp)
				break
			}
			tmp.Prerelease = Prerelease(string(prerelease))
		}
		if s == '+' {
			if metadata == "" {
				// There was no '-' found, so the code for '-' above did not run
				metadata = Metadata(ver[i+1:])
			}
			if re.MatchString(string(metadata)) {
				msg = fmt.Sprintf("build metadata in '%s' is not valid semver", ver)
				hlp = fmt.Sprintf("build metadata %s#spec-item-10", sharedHelp)
				break
			}
			tmp.Metadata = metadata
		}

	}
	if msg != "" {
		sts = status.Fail(&status.Args{
			Message: msg,
			Help:    hlp,
		})
	}
	if status.IsSuccess(sts) {
		*me = tmp
	}
	return sts
}

func (me *Version) captureMMP(ver types.Version, pos int, buf []byte) (newpos int, msg string) {
	switch pos {
	case 0:
		me.Major = Digits(buf)
	case 1:
		if me.Major == "" {
			msg = fmt.Sprintf("version '%s' contains minor version but no major version", ver)
			break
		}
		me.Minor = Digits(buf)
	case 2:
		if me.Minor == "" {
			msg = fmt.Sprintf("version '%s' contains patch but no minor version", ver)
			break
		}
		if me.Major == "" {
			msg = fmt.Sprintf("version '%s' contains patch but no major version", ver)
			break
		}
		me.Patch = Digits(buf)
	}
	newpos = pos + 1
	return newpos, msg
}

func (me *Version) String() string {
	var ver string
	for i := verRelease; i >= verMajor; i-- {
		switch i {
		case verRelease:
			if me.Revision != "" {
				ver = fmt.Sprintf("~r%s", me.Revision)
			}
		case verMetadata:
			if me.Metadata != "" {
				ver = fmt.Sprintf("+%s%s", me.Metadata, ver)
			}
		case verPrerelease:
			if me.Prerelease != "" {
				ver = fmt.Sprintf("-%s%s", me.Prerelease, ver)
			}
		case verPatch:
			if me.Patch != "" {
				ver = fmt.Sprintf(".%s%s", me.Patch, ver)
			}
		case verMinor:
			if me.Minor != "" {
				ver = fmt.Sprintf(".%s%s", me.Minor, ver)
			}
		case verMajor:
			if me.Major != "" {
				ver = fmt.Sprintf("%s%s", me.Major, ver)
			}
		}
	}
	return ver
}

func (me *Version) GetRaw() types.Version {
	return me.raw
}

func (me *Version) GetMajor() Digits {
	me.checkParsed("Major")
	return me.Major
}

func (me *Version) GetMinor() Digits {
	me.checkParsed("Minor")
	return me.Minor
}

func (me *Version) GetPatch() Digits {
	me.checkParsed("Patch")
	return me.Patch
}

func (me *Version) GetPrerelease() Prerelease {
	me.checkParsed("Prerelease")
	return me.Prerelease
}

func (me *Version) GetMetadata() Metadata {
	me.checkParsed("Metadata")
	return me.Metadata
}

func (me *Version) GetMajorMinor() MajorMinor {
	me.checkParsed("MajorMinor")
	var mm MajorMinor
	if me.Minor == "" {
		mm = MajorMinor(me.Major)
	} else {
		mm = MajorMinor(fmt.Sprintf("%s.%s", me.Major, me.Minor))
	}
	return mm
}

func (me *Version) GetIdentifier() types.Version {
	me.checkParsed("Version")
	return types.Version(me.String())
}

func (me *Version) GetRevision() Revision {
	me.checkParsed("Revision")
	return me.Revision
}

func (me *Version) checkParsed(f string) {
	if me.raw == "" {
		panic(fmt.Sprintf("accessing Version.Get%s() before initializing with Version.Parse()", f))
	}
}
