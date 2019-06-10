package messages

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)


type Message struct {
	Time MessageTime
	Source  MessageAddress
	Topic MessageTopic
	Text MessageText

	//PayLoad
}


type PayLoad struct {
	Topic MessageTopic
	Text MessageText
}


const (
	Package                = "messages"
	InterfaceTypeMessage   = "*" + Package + ".Message"
	InterfaceTypeSubTopic  = "*" + Package + ".SubTopic"
	InterfaceTypeSubTopics = "*" + Package + ".SubTopics"
	InterfaceTypeError     = "error"
)


func (me MessageText) ToMessageAddress() MessageAddress {

	return MessageAddress(me.String())
}


func (me MessageAddress) ConstructMessage(to MessageAddress, subtopic SubTopic, text MessageText) Message {

	//var err error
	var msgTemplate Message

	msgTemplate = Message{
		Source: me,
		Topic: MessageTopic{
			Address: to,
			SubTopic: subtopic,
		},
		Text: text,
	}

	return msgTemplate
}


func (me MessageAddress) ConstructTopic(to MessageAddress, subtopic SubTopic) MessageTopic {

	//var err error
	var topicTemplate MessageTopic

	topicTemplate = MessageTopic{
		Address:  to,
		SubTopic: subtopic,
	}

	return topicTemplate
}


func (me *MessageAddress) ProduceError(msg string, a ...interface{}) error {

	if me == nil {
		return errors.New(fmt.Sprintf(msg, a...))
	} else {
		return errors.New(fmt.Sprintf(me.String() + ": " + msg, a...))
	}
}


func ProduceError(me MessageAddress, msg string, a ...interface{}) error {

	if me == "" {
		return errors.New(fmt.Sprintf(msg, a...))
	} else {
		return errors.New(fmt.Sprintf(me.String() + ": " + msg, a...))
	}
}


type MessageTime time.Time
func (me *MessageTime) IsNil() bool {

	if *me == MessageTime(DefaultNilTime) {
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
func (me *MessageTime) Unix() int64 {

	return time.Time(*me).Unix()
}


type Uuid uuid.UUID
func (me *Uuid) IsNil() bool {

	if me == nil {
		return true
	}

	return false
}
func (me *Uuid) String() string {
	return uuid.UUID(*me).String()
}
func (me *Uuid) ToMessageType() MessageText {
	// return uuid.UUID(*me).String()
	return MessageText(me.String())
}
func (me *Uuid) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("message address uuid is nil")
	}

	return err
}
func EnsureUuidNotNil(me *Uuid) error {

	var err error

	switch {
		case me == nil:
			err = errors.New("message address uuid is nil")
	}

	return err
}


type MessageAddress string
type MessageAddresses []MessageAddress
func GenerateAddress() MessageAddress {

	return MessageAddress(uuid.New().String())
}
func (me *MessageAddress) String() string {

	return string(*me)
}
func (me *MessageAddress) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("message address is nil")

		case *me == "":
			err = errors.New("message address is empty")
	}

	return err
}


type MessageText string
func (me *MessageText) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("message text is nil")

		case *me == "":
			err = errors.New("message text is empty")
	}

	return err
}
func (me *MessageText) String() string {

	return string(*me)
}
func (me *MessageText) ByteArray() []byte {

	return []byte(*me)
}

