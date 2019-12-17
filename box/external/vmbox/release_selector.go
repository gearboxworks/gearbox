package vmbox

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"github.com/gearboxworks/go-status/only"
	"time"
)

type ReleaseSelector struct {
	// These are considered to be AND-ed together.
	FromDate        time.Time
	UntilDate       time.Time
	SpecificVersion string
	RegexpVersion   string
	Latest          *bool
}

/*
Updates the following:
   me.VmIsoVersion    string
   me.VmIsoFile       string
   me.VmIsoUrl 		string
   me.VmIsoRelease    Release
*/
func (me *Releases) SelectRelease(selector ReleaseSelector) (*Release, error) {

	var err error
	var r *Release

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//err = me.UpdateReleases()
		//if err != nil {
		//	break
		//}

		// For now just select the latest.
		me.Selected = me.Latest
		r = me.Selected

		eblog.Debug(entity.VmBoxEntityName, "selecting the latest release == %s", me.Latest.Version)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return r, err
}
