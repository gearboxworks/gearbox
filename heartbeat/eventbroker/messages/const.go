package messages

import (
	"time"
)

const (
	TopicSeparator = "/%s/%s"
	BroadcastAddress = "broadcast"
)


var DefaultNilTime = time.Time{}
var IsEmptySubTopic SubTopic

// 	DefaultExitString = "exit"
//	SubTopicState     = messages.SubTopic("state")
