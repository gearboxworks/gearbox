package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"regexp"
	"strconv"
	"strings"
)

type DottedVersion struct {
	raw        string
	Major      string `json:"major,omitempty"`
	Minor      string `json:"minor,omitempty"`
	Patch      string `json:"patch,omitempty"`
	Prerelease string `json:"prerelease,omitempty"`
	Metadata   string `json:"metadata,omitempty"`
	Revision   string `json:"revision,omitempty"`
	Error      error  `json:"error,omitempty"`
}

func NewDottedVersion() *DottedVersion {
	return &DottedVersion{}
}

func (me *DottedVersion) Parse(ver string) (err error) {
	newBuf := true
	parts := strings.Split(ver, "~")
	tmp := DottedVersion{raw: ver}
	re := regexp.MustCompile("[^A-Za-z0-9.-]")
	for range only.Once {
		if len(parts) >= 2 {
			if len(parts[1]) == 0 {
				err = fmt.Errorf("revision following '~' in '%s' is empty", ver)
				break
			}
			if parts[1][0] != 'r' {
				err = fmt.Errorf("revision following '~' in '%s' must begin with 'r'", ver)
				break
			}
			tmp.Revision = parts[1][1:]
			_, err = strconv.Atoi(tmp.Revision)
			if err != nil {
				err = fmt.Errorf("revision following '~r' in '%s' must be and integer", ver)
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
				pos, err = tmp.captureMMP(ver, pos, buf)
				if err != nil {
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
				err = fmt.Errorf("non-integer %s version in '%s'",
					[]string{"major", "minor", "patch"}[pos],
					ver,
				)
				break
			}
			if err != nil || done {
				break
			}
		}
		if err != nil {
			break
		}
		_, err = tmp.captureMMP(ver, pos, buf)
		if err != nil {
			break
		}
		metadata := ""
		sharedHelp := "can only contain letters, numbers, periods or dashes. See the SemVer docs for more: https://semver.org/"
		if s == '-' {
			prerelease := ver[i+1:]
			idx := strings.Index(prerelease, "+")
			if idx != -1 {
				metadata = prerelease[idx+1:]
				prerelease = prerelease[0:idx]
				s = '+'
			}
			if re.MatchString(prerelease) {
				err = util.AddHelpToError(
					fmt.Errorf("pre-release in '%s' is invalid semver", ver),
					fmt.Sprintf("pre-release %s#spec-item-10", sharedHelp),
				)
				break
			}
			tmp.Prerelease = string(prerelease)
		}
		if s == '+' {
			if metadata == "" {
				// There was no '-' found, so the code for '-' above did not run
				metadata = ver[i+1:]
			}
			if re.MatchString(metadata) {
				err = util.AddHelpToError(
					fmt.Errorf("build metadata in '%s' is not valid semver", ver),
					fmt.Sprintf("build metadata %s#spec-item-10", sharedHelp),
				)
				break
			}
			tmp.Metadata = metadata
		}

	}
	if err != nil {
		me.Error = err
	} else {
		*me = tmp
	}
	return err
}

func (me *DottedVersion) captureMMP(ver string, pos int, buf []byte) (int, error) {
	var err error
	switch pos {
	case 0:
		me.Major = string(buf)
	case 1:
		if me.Major == "" {
			err = fmt.Errorf("version '%s' contains minor version but no major version", ver)
			break
		}
		me.Minor = string(buf)
	case 2:
		if me.Minor == "" {
			err = fmt.Errorf("version '%s' contains patch but no minor version", ver)
			break
		}
		if me.Major == "" {
			err = fmt.Errorf("version '%s' contains patch but no major version", ver)
			break
		}
		me.Patch = string(buf)
	}
	pos++
	return pos, err
}

const (
	dvMajor = iota
	dvMinor
	dvPatch
	dvPrerelease
	dvMetadata
	dvRelease
)

func (me *DottedVersion) String() string {
	var ver string
	for i := dvRelease; i >= dvMajor; i-- {
		switch i {
		case dvRelease:
			if me.Revision != "" {
				ver = fmt.Sprintf("~r%s", me.Revision)
			}
		case dvMetadata:
			if me.Metadata != "" {
				ver = fmt.Sprintf("+%s%s", me.Metadata, ver)
			}
		case dvPrerelease:
			if me.Prerelease != "" {
				ver = fmt.Sprintf("-%s%s", me.Prerelease, ver)
			}
		case dvPatch:
			if me.Patch != "" {
				ver = fmt.Sprintf(".%s%s", me.Patch, ver)
			}
		case dvMinor:
			if me.Minor != "" {
				ver = fmt.Sprintf(".%s%s", me.Minor, ver)
			}
		case dvMajor:
			if me.Major != "" {
				ver = fmt.Sprintf("%s%s", me.Major, ver)
			}
		}
	}
	return ver
}

func (me *DottedVersion) GetRaw() string {
	return me.raw
}

func (me *DottedVersion) GetMajor() string {
	me.checkParsed("Major")
	return me.Major
}

func (me *DottedVersion) GetMinor() string {
	me.checkParsed("Minor")
	return me.Minor
}

func (me *DottedVersion) GetPatch() string {
	me.checkParsed("Patch")
	return me.Patch
}

func (me *DottedVersion) GetPrerelease() string {
	me.checkParsed("Prerelease")
	return me.Prerelease
}

func (me *DottedVersion) GetMetadata() string {
	me.checkParsed("Metadata")
	return me.Metadata
}

func (me *DottedVersion) GetMajorMinor() string {
	me.checkParsed("MajorMinor")
	var mm string
	if me.Minor == "" {
		mm = me.Major
	} else {
		mm = fmt.Sprintf("%s.%s", me.Major, me.Minor)
	}
	return mm
}

func (me *DottedVersion) GetVersion() string {
	me.checkParsed("Version")
	return me.String()
}

func (me *DottedVersion) GetRevision() string {
	me.checkParsed("Revision")
	return me.Revision
}

func (me *DottedVersion) checkParsed(f string) {
	if me.raw == "" {
		panic(fmt.Sprintf("accessing DottedVersion.Get%s() before initializing with DottedVersion.Parse()", f))
	}
}
