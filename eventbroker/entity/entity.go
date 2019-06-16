package entity

import "gearbox/eventbroker/messages"

const (
	ChannelEntityName = "eventbroker-channels"
	DaemonEntityName = "eventbroker-daemon"
	MqttClientEntityName = "eventbroker-mqttclient"
	NetworkEntityName = "eventbroker-network"
	BroadcastEntityName = "broadcast"
	SelfEntityName = "self"

	UnfsdEntityName = "unfsd"
	MqttBrokerEntityName = "mqtt"
	ApiEntityName = "api"
	VmBoxEntityName = "vm"
)
//var SelfEntityName = "self"

var AllEntities = messages.MessageAddresses{ChannelEntityName, NetworkEntityName, DaemonEntityName, MqttClientEntityName}
var PartialEntities = messages.MessageAddresses{NetworkEntityName, DaemonEntityName, MqttClientEntityName}
