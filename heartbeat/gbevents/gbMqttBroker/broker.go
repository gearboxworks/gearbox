package gbMqttBroker

import (
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"net/url"
	"time"
)


func (me *Mqtt) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

	var _args Args
	var sts status.Status

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy MQTT broker config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}


		// Setup Broker.
		// ConfigureConfig is normally used with args, but good enough to create a default config.
		_args.Broker.EntityId = DefaultBrokerEntityId
		// _args.Broker.State = false
		_args.Broker.RestartAttempts = DefaultRetries
		_args.Broker.config, err = broker.ConfigureConfig([]string{""})
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to parse MQTT broker config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}
		_args.Broker.config = broker.DefaultConfig

		if _args.ServerURL == nil {
			_args.ServerURL, _ = url.Parse(DefaultServerUrl)
		}

		if _args.BrokerWorkers == 0 {
			_args.BrokerWorkers = DefaultWorkers
		}

		_args.Broker.instance, err = broker.NewBroker(_args.Broker.config)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("error starting MQTT broker").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}


		// Setup client.
		_args.Client.EntityId = DefaultClientEntityId
		// _args.Client.State = false
		_args.Client.RestartAttempts = DefaultRetries
		_args.Client.config = mqtt.NewClientOptions()
		_args.Client.config.AddBroker(_args.ServerURL.String())
		//_args.Client.config.SetUsername(me.Server.User.Username())
		//password, _ := me.Server.User.Password()
		//_args.Client.config.SetPassword(password)
		_args.Client.config.SetClientID(_args.EntityId.String())
		_args.Client.instance = mqtt.NewClient(_args.Client.config)

		//		topic := messages.CreateTopic(_args.EntityId.String())
		//		_args.Config.SetWill(topic, "Last will and testament", topics.QosFailure, true)

		// _args.Server = fmt.Sprintf("tcp://%s:%s/", _args.Config.Host, _args.Config.Port)

		*me = Mqtt(_args)
		logger.Debug("MQTT init %s.", me.EntityId.String())
	}

	return sts
}


func (me *Mqtt) StartHandler() status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		go func() {
			me.StartBrokerHandler().Log()
		}()

		//go func() {
		//	me.StartClientHandler().Log()
		//}()

		logger.Debug("MQTT started OK on %s.", me.EntityId.String())
		sts = status.Success("MQTT started OK on %s", me.EntityId.String())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	return sts
}


func (me *Mqtt) StartBrokerHandler() status.Status {

	var sts status.Status
	var err error

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		me.Broker.Task, err = tasks.StartTask(initMqttBroker, startMqttBroker, monitorMqttBroker, stopMqttBroker, me)
		if err != nil {
			sts = me.Broker.Sts
			break
		}

		//time.Sleep(time.Second * 10)
		//me.Broker.Task.Stop()


		time.Sleep(time.Hour * 24)

		/*
		for me.Broker.Task.RunCounter = 0; me.Broker.Task.RunCounter < me.Broker.RestartAttempts; me.Broker.Task.RunCounter++ {

			if me.Broker.Task.RunCounter == 0 {
				logger.Debug("MQTT broker %s started.", me.EntityId.String())
			} else {
				//me.Broker.State = false
				logger.Debug("MQTT broker %s restart attempt %d.", me.EntityId.String(), me.Broker.Task.RunCounter)
			}

			// .			me.Broker.instance.Start()

			if me.Broker.Task.RunCounter == 0 {
				//me.Broker.State = true
				logger.Debug("MQTT broker %s started.", me.EntityId.String())
			} else {
				//me.Broker.State = false
				logger.Debug("MQTT broker %s restart attempt %d.", me.EntityId.String(), me.Broker.Task.RunCounter)
			}

			logger.Debug("MQTT broker stopped %s.", me.EntityId.String())
			sts = status.Success("MQTT broker exited with signal %v.")
		}
		*/
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Broker.Sts = sts

	return sts
}


func initMqttBroker(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		//me := (i[0]).(*Mqtt)
		me := i[0].(*Mqtt)

		sts := me.EnsureNotNil()
		if is.Error(sts) {
			err = sts
			break
		}

		//task.retryLimit = DefaultRetries
		//task.retryDelay = time.Second * 5

		me.Broker.Sts = status.Success("MQTT broker %s initialized OK", me.ServerURL.String())
		me.Broker.Sts.Log()

		err = nil
	}

	return err
}


func startMqttBroker(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		me := i[0].(*Mqtt)

		sts := me.EnsureNotNil()
		if is.Error(sts) {
			err = sts
			break
		}

		if task.GetRetryCounter() == 0 {
			me.Broker.Sts = status.Success("MQTT broker %s started", me.ServerURL.String())
		} else {
			me.Broker.Sts = status.Success("MQTT broker %s restart attempt %d", me.ServerURL.String(), task.GetRetryCounter())
		}
		me.Broker.Sts.Log()

		// This function will terminate almost immediately with no indication of state as it contains GoRoutines.
		// statusMqttBroker() will perform the actual status check.
		me.Broker.instance.Start()

		me.Broker.Sts = status.Success("MQTT broker %s init OK", me.EntityId.String())
		me.Broker.Sts.Log()

		err = nil
	}

	return err
}


func monitorMqttBroker(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		me := i[0].(*Mqtt)

		sts := me.EnsureNotNil()
		if is.Error(sts) {
			err = sts
			break
		}

		me.Broker.Sts = status.Success("checking MQTT broker %s status", me.ServerURL.String())
		me.Broker.Sts.Log()


/*
		// So, let's just check on things in a loop.
		localClient := mqtt.NewClient(me.Client.config)
		if localClient == nil {
			me.Broker.Sts = status.Fail().
				SetMessage("MQTT broker failed to %s start", me.ServerURL.String()).
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			err = me.Broker.Sts
			break
		}
		defer localClient.Disconnect(0)

		localToken := localClient.Connect()
		if localToken == nil {
			me.Broker.Sts = status.Fail().
				SetMessage("MQTT broker failed to %s start", me.ServerURL.String()).
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			err = me.Broker.Sts
			break
		}

		for !localToken.WaitTimeout(3 * time.Second) {
		}

		if err := localToken.Error(); err != nil {
			me.Broker.Sts = status.Wrap(err).
				SetMessage("MQTT broker failed to %s start", me.ServerURL.String()).
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			err = me.Broker.Sts
			break
		}
*/

		err = nil
	}

	return err
}


func stopMqttBroker(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		me := i[0].(*Mqtt)

		sts := me.EnsureNotNil()
		if is.Error(sts) {
			err = sts
			break
		}

		me.Broker.Sts = status.Success("MQTT broker %s stopped", me.ServerURL.String())
		me.Broker.Sts.Log()

		err = nil
	}

	return err
}




func (me *Mqtt) StartClientHandler() status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if me.Client.instance == nil {
			sts = status.Fail().
				SetMessage("MQTT client not configured").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		//me.Broker.State = true
		logger.Debug("MQTT client started %s.", me.EntityId.String())

		sts = me.clientConnect()
		if is.Error(sts) {
			break
		}

		//sts = me.subscribe(topic, nil)
		if is.Error(sts) {
			break
		}

		s := daemon.WaitForSignal()

		//me.Broker.State = false
		logger.Debug("MQTT client stopped %s.", me.EntityId.String())
		sts = status.Success("MQTT client exited with signal %v.", s)
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Client.Sts = sts

	return sts
}


func (me *Mqtt) clientConnect() (status.Status) {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		me.Client.token = me.Client.instance.Connect()
		if me.Client.token == nil {
			sts = status.Fail().
				SetMessage("unexpected software error").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		for !me.Client.token.WaitTimeout(3 * time.Second) {
		}

		if err := me.Client.token.Error(); err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to connect to MQTT broker %s", me.ServerURL.String()).
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		sts = status.Success("MQTT client connected to %s OK.", me.ServerURL.String())
	}

	return sts
}


func (me *Mqtt) subscribe(topic messages.Topic, foo mqtt.MessageHandler) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		cb := msgCallback{Topic: messages.Topic(topic), Function: foo2}
		me.Client.instance.Subscribe(cb.Topic.String(), 0, cb.Function)

		sts = status.Success("MQTT client subscribed OK")
	}

	return sts
}


func foo2(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


func callbackFunc(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}

