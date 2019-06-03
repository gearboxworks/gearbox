package messages

import (
	"time"
)


type Message struct {
	Time MessageTime
	Source  MessageAddress
	// Destination  MessageAddress
	Topic Topic
	Text MessageText
}

type MessageTime time.Time
func (me *MessageTime) IsNil() bool {

	if *me == MessageTime(defaultNilTime) {
		return true
	}

	return false
}
func (me *MessageTime) Now() MessageTime {

	return MessageTime(time.Now())
}
func (me *MessageTime) Convert() time.Time {

	return time.Time(*me)
}

type MessageAddress string
func (me *MessageAddress) String() string {

	return string(*me)
}
func (me *MessageAddress) IsNil() bool {

	// 	var emptyMA MessageAddress
	//sts = status.Warn("").
	//	SetMessage("client src address is empty").
	//	SetAdditional("", ).
	//	SetData("").
	//	SetHelp(status.AllHelp, help.ContactSupportHelp())

	if me == nil {
		return true
	}

	return false
}


type MessageText string
func (me *MessageText) String() string {

	return string(*me)
}
func (me *MessageText) ByteArray() []byte {

	return []byte(*me)
}

var defaultNilTime = time.Time{}
