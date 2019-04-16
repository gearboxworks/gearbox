package test

import (
	"gearbox/gearspec"
	"gearbox/status/is"
	"gearbox/types"
	"testing"
)

type wants struct {
	StackId     types.StackId
	Identifier  gearspec.Identifier
	Persistable gearspec.Identifier
	Expanded    gearspec.Identifier
}

var results = map[gearspec.Identifier]wants{
	"gears.gearbox.works/wordpress/dbserver": wants{
		StackId:     "gears.gearbox.works/wordpress",
		Identifier:  "gears.gearbox.works/wordpress/dbserver",
		Persistable: "wordpress/dbserver",
		Expanded:    "gears.gearbox.works/wordpress/dbserver",
	},
	"wordpress/dbserver": wants{
		StackId:     "gears.gearbox.works/wordpress",
		Identifier:  "gears.gearbox.works/wordpress/dbserver",
		Persistable: "wordpress/dbserver",
		Expanded:    "gears.gearbox.works/wordpress/dbserver",
	},
}

func TestIdentifier(t *testing.T) {

	t.Run("GetNamedStackId", func(t *testing.T) {
		for gearspecid, wants := range results {
			nsid, sts := gearspecid.GetNamedStackId()
			if is.Error(sts) {
				t.Errorf("For Gearspec ID '%s': %s",
					gearspecid,
					sts.Message(),
				)
				continue
			}
			if nsid != wants.StackId {
				t.Errorf("For Gearspec ID '%s' wanted '%s' got '%s'",
					gearspecid,
					wants.StackId,
					nsid,
				)
			}
		}
	})

	t.Run("GetIdentifier", func(t *testing.T) {
		for gearspecid, wants := range results {
			_gsid, sts := gearspecid.GetIdentifier()
			if is.Error(sts) {
				t.Errorf("For Gearspec ID '%s': %s",
					gearspecid,
					sts.Message(),
				)
				continue
			}
			if _gsid != wants.Identifier {
				t.Errorf("For Gearspec ID '%s' wanted '%s' got '%s'",
					gearspecid,
					wants.Identifier,
					_gsid,
				)
			}
		}
	})

	t.Run("GetPersistableIdentifier", func(t *testing.T) {
		for gearspecid, wants := range results {
			_gsid, sts := gearspecid.GetPersistableIdentifier()
			if is.Error(sts) {
				t.Errorf("For Gearspec ID '%s': %s",
					gearspecid,
					sts.Message(),
				)
				continue
			}
			if _gsid != wants.Persistable {
				t.Errorf("For Gearspec ID '%s' wanted '%s' got '%s'",
					gearspecid,
					wants.Persistable,
					_gsid,
				)
			}
		}
	})

	t.Run("GetExpandedIdentifier", func(t *testing.T) {
		for gearspecid, wants := range results {
			_gsid, sts := gearspecid.GetExpandedIdentifier()
			if is.Error(sts) {
				t.Errorf("For Gearspec ID '%s': %s",
					gearspecid,
					sts.Message(),
				)
				continue
			}
			if _gsid != wants.Expanded {
				t.Errorf("For Gearspec ID '%s' wanted '%s' got '%s'",
					gearspecid,
					wants.Expanded,
					_gsid,
				)
			}
		}
	})

	return
}
