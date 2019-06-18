package entity

import "gearbox/eventbroker/messages"

const (
	ChannelEntityName    = "eventbroker-channels"
	DaemonEntityName     = "eventbroker-daemon"
	MqttClientEntityName = "eventbroker-mqttclient"
	NetworkEntityName    = "eventbroker-network"
	BroadcastEntityName  = "broadcast"
	SelfEntityName       = "self"
	UnfsdEntityName      = "unfsd"
	MqttBrokerEntityName = "mqtt"

	HeartbeatEntityName  = "heartbeat"
	ApiEntityName        = "api"
	VmBoxEntityName      = "vm"
	VmUpdateEntityName   = "update"
	VmEntityName         = "Gearbox"
)

//var SelfEntityName = "self"

var AllEntities = messages.MessageAddresses{ChannelEntityName, NetworkEntityName, DaemonEntityName, MqttClientEntityName}
var PartialEntities = messages.MessageAddresses{NetworkEntityName, DaemonEntityName, MqttClientEntityName}
