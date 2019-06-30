package entity

import "gearbox/eventbroker/msgs"

const (
	ChannelEntityName    = "eventbroker-channels"
	DaemonEntityName     = "eventbroker-daemon"
	MqttClientEntityName = "eventbroker-mqttclient"
	NetworkEntityName    = "eventbroker-network"
	BroadcastEntityName  = "broadcast"
	SelfEntityName       = "self"
	UnfsdEntityName      = "unfsd"
	MqttBrokerEntityName = "mqtt"

	HeartbeatEntityName = "heartbeat"
	ApiEntityName       = "api"
	VmBoxEntityName     = "vm"
	VmUpdateEntityName  = "update"
	VmEntityName        = "Gearbox"
)

//var SelfEntityName = "self"

var AllEntities = msgs.Addresses{ChannelEntityName, NetworkEntityName, DaemonEntityName, MqttClientEntityName}
var PartialEntities = msgs.Addresses{NetworkEntityName, DaemonEntityName, MqttClientEntityName}
