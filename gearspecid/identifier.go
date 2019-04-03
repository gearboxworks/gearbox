package gsid

import (
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

//
// Gearspec Example:
//
// 		gearbox.works/wordpress/dbserver:1
//

type Identifiers []Identifier
type Identifier string

func (me Identifier) GetNamedStackId() (sid types.StackId, sts status.Status) {
	for range only.Once {
		gsi := &Id{}
		sts = gsi.Parse(me)
		if is.Error(sts) {
			break
		}
		sid = gsi.GetStackId()
	}
	return sid, sts
}

func (me Identifier) GetIdentifier() (gs Identifier, sts status.Status) {
	gsi := Id{}
	sts = gsi.Parse(me)
	if is.Success(sts) && gsi.Authority == "" {
		gsi.Authority = global.DefaultAuthority
	}
	return Identifier(gsi.String()), sts
}

//
// StackRole.GetPersistableIdentifier()
//
// Returns a Gearspec without authority if authority is "gearbox.works"
//
// Used to write values to the gearbox.json configuration file
// to keep things simple for the user/reader.
//
func (me Identifier) GetPersistableIdentifier() (gs Identifier, sts status.Status) {
	gsi := Id{}
	sts = gsi.Parse(me)
	if is.Success(sts) && gsi.Authority == global.DefaultAuthority {
		gsi.Authority = ""
	}
	return Identifier(gsi.String()), sts
}

func (me Identifier) GetExpandedIdentifier() (gs Identifier, sts status.Status) {
	for range only.Once {
		gs = me
		gsi := NewGearspecId()
		sts = gsi.Parse(me)
		if is.Error(sts) {
			break
		}
		gs = Identifier(gsi.String())
	}
	return gs, sts
}
