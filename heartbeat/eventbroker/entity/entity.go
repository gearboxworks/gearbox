package entity

import "gearbox/heartbeat/eventbroker/messages"

const (
	ChannelEntityName = "eventbroker-channels"
	DaemonEntityName = "eventbroker-daemon"
	MqttClientEntityName = "eventbroker-mqttclient"
	NetworkEntityName = "eventbroker-network"
	BroadcastEntityName = "broadcast"
	SelfEntityName = "self"
)
//var SelfEntityName = "self"

var AllEntities = messages.MessageAddresses{ChannelEntityName, NetworkEntityName, DaemonEntityName, MqttClientEntityName}
var PartialEntities = messages.MessageAddresses{NetworkEntityName, DaemonEntityName, MqttClientEntityName}
