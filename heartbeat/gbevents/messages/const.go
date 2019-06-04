package messages

import "time"

const (
	SubTopicState  = SubTopic("state")
	SubTopicStart  = SubTopic("start")
	SubTopicStop   = SubTopic("stop")
	SubTopicStatus = SubTopic("state")
	SubTopicError  = SubTopic("error")
	SubTopicGlob   = SubTopic("*")

	SubTopicUnregister = SubTopic("unregister")
	SubTopicRegister = SubTopic("register")
	SubTopicSubscribe = SubTopic("subscribe")
	SubTopicUnsubscribe = SubTopic("unsubscribe")

	MessageStateIdle = MessageText("idle")
	MessageStateUnconfigured = MessageText("unconfigured")
	MessageStateError = MessageText("error")
	MessageStateUp = MessageText("up")
	MessageStateDown = MessageText("down")
	MessageStateStarting = MessageText("starting")
	MessageStateStopping = MessageText("stopping")

	TopicSeparator = "/%s/%s"
)

var DefaultNilTime = time.Time{}
var IsEmptySubTopic SubTopic

// 	DefaultExitString = "exit"
//	SubTopicState     = messages.SubTopic("state")