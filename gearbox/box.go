package gearbox

import (
	"gearbox/box"
	"gearbox/only"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
)

func newBox(me Gearboxer, args box.Args) (bx *box.Box, sts status.Status) {

	bx, _ = box.New(me.GetOsBridge(), args)
	//sts = bx.Initialize()

	return bx, sts
}


func (me *Gearbox) BoxDaemon(args box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		bx, sts = box.New(me.GetOsBridge(), args)
		if is.Error(sts) {
			break
		}

		sts = bx.BoxDaemon()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) StartBox(args box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		bx, sts = box.New(me.GetOsBridge(), args)
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

func (me *Gearbox) StopBox(args box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		bx, sts = box.New(me.GetOsBridge(), args)
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

func (me *Gearbox) RestartBox(args box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		bx, sts = box.New(me.GetOsBridge(), args)
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

func (me *Gearbox) PrintBoxStatus(args box.Args) (sts status.Status) {

	for range only.Once {
		var bx *box.Box

		bx, sts = box.New(me.GetOsBridge(), args)
		if is.Error(sts) {
			break
		}

		sts := bx.GetState()
		if is.Error(sts) {
			break
		}

		//		var state string
		//		state, sts = sts.GetString()
		//		meaning := box.GetStateMeaning(box.State(state))
		//		if meaning == "" {
		//			fmt.Println(box.GetStateMeaning(box.UnknownState))
		//			break
		//		}
		//		fmt.Println(meaning)
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) CreateBox(box.Args) status.Status {
	panic("implement me")
}




//func newBox(me Gearboxer, args box.Args) (bx *box.Box, sts status.Status) {
//
//	bx = box.NewBox(me.GetOsBridge(), args)
//	sts = bx.Initialize()
//
//	return bx, sts
//}
//
//func (me *Gearbox) StartBox(args box.Args) (sts status.Status) {
//
//	for range only.Once {
//
//		var bx *box.Box
//
//		bx, sts = newBox(me, args)
//		if is.Error(sts) {
//			break
//		}
//
//		sts = bx.Start()
//		if is.Error(sts) {
//			break
//		}
//	}
//
//	return sts
//}
//
//func (me *Gearbox) StopBox(args box.Args) (sts status.Status) {
//
//	for range only.Once {
//
//		var bx *box.Box
//
//		bx, sts = newBox(me, args)
//		if is.Error(sts) {
//			break
//		}
//
//		sts = bx.Stop()
//		if is.Error(sts) {
//			break
//		}
//	}
//
//	return sts
//}
//
//func (me *Gearbox) RestartBox(args box.Args) (sts status.Status) {
//
//	for range only.Once {
//
//		var bx *box.Box
//
//		bx, sts = newBox(me, args)
//		if is.Error(sts) {
//			break
//		}
//
//		sts = bx.Restart()
//		if is.Error(sts) {
//			break
//		}
//	}
//
//	return sts
//}
//
//func (me *Gearbox) PrintBoxStatus(args box.Args) (sts status.Status) {
//	for range only.Once {
//
//		var bx *box.Box
//
//		bx, sts = newBox(me, args)
//		if is.Error(sts) {
//			break
//		}
//
//		_, sts := bx.GetState()
//		if is.Error(sts) {
//			break
//		}
//	}
//
//	return sts
//}
//
//func (me *Gearbox) CreateBox(args box.Args) (sts status.Status) {
//
//	for range only.Once {
//
//		var bx *box.Box
//
//		bx, sts = newBox(me, args)
//		if is.Error(sts) {
//			break
//		}
//
//		_, sts = bx.CreateBox()
//		if is.Error(sts) {
//			break
//		}
//
//		sts = status.Success("%s VM created", global.Brandname)
//	}
//	return sts
//
//}
