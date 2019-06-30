package gearbox

import (
	"gearbox/box"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

func newBox(me Gearboxer, args *box.Args) (bx *box.Box, sts status.Status) {

	args.SetOsBridge(me.GetOsBridge())
	bx, _ = box.New(args)
	//sts = bx.Initialize()

	return bx, sts
}

func (me *Gearbox) BoxDaemon(args *box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		args.SetOsBridge(me.GetOsBridge())
		bx, sts = box.New(args)
		if is.Error(sts) {
			break
		}

		sts = bx.RunAsDaemon()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) StartBox(args *box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		args.SetOsBridge(me.GetOsBridge())
		bx, sts = box.New(args)
		if is.Error(sts) {
			break
		}

		sts = bx.StartBox()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) StopBox(args *box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		args.SetOsBridge(me.GetOsBridge())
		bx, sts = box.New(args)
		if is.Error(sts) {
			break
		}

		sts = bx.StopBox()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) RestartBox(args *box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		args.SetOsBridge(me.GetOsBridge())
		bx, sts = box.New(args)
		if is.Error(sts) {
			break
		}

		sts = bx.RestartBox()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) PrintBoxStatus(args *box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		args.SetOsBridge(me.GetOsBridge())
		bx, sts = box.New(args)
		if is.Error(sts) {
			break
		}

		sts := bx.GetState()
		if is.Error(sts) {
			break
		}

	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) CreateBox(*box.Args) status.Status {
	panic("implement me")
}
