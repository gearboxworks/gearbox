package external

//func New(OsBridge osbridge.OsBridger, args ...Args) (*EventBroker, error) {
//
//	var _args Args
//	var err error
//
//	me := &EventBroker{}
//
//	for range only.Once {
//
//		if len(args) > 0 {
//			_args = args[0]
//		}
//
//		if _args.Boxname == "" {
//			_args.Boxname = global.Brandname
//		}
//
//		if _args.EntityId == "" {
//			_args.EntityId = DefaultEntityName
//		}
//
//		_args.osSupport = OsBridge
//		foo := box.Args{}
//		err = copier.Copy(&foo, &_args)
//		if err != nil {
//			err = msgs.MakeError(me.EntityId,"unable to copy config args")
//			break
//		}
//
//		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.osSupport.GetAdminRootDir(), defaultPidFile))
//
//		*me = EventBroker(_args)
//
//
//		me.State.SetWant(states.StateIdle)
//		if me.State.SetNewState(states.StateIdle, err) {
//			eblog.Debug(me.EntityId, "init complete")
//		}
//	}
//
//	//channels.PublishCallerState(&me.Channels, &me.EntityId, &me.State)
//	//eblog.LogIfNil(me, err)
//	//eblog.LogIfError(me.EntityId, err)
//
//	return me, err
//}
