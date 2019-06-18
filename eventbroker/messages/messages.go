package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/eventbroker/only"
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
)


func (me MessageText) ToMessageAddress() MessageAddress {

	return MessageAddress(me.String())
}


func (me *Message) String() string {

	return fmt.Sprintf(`Time:%d  Source:%s  Topic:%s  Text:%s`,
		me.Time.Unix(),
		me.Source.String(),
		me.Topic.String(),
		me.Text.String(),
	)
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

	if me.IsNil() {
		return ""
	}

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
func GenerateAddress() *MessageAddress {

	u := MessageAddress(uuid.New().String())

	return &u
}
func (me *MessageAddress) String() string {

	if me.EnsureNotNil() != nil {
		return ""
	}

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

func (me *Message) Validate() error {

	var err error

	for range only.Once {
		//err = me.Text.EnsureNotNil()
		//if err != nil {
		//	break
		//}

		err = me.Topic.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.Source.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return err
}

func (me *Message) ToMessageText() MessageText {

	var err error
	var j []byte

	for range only.Once {
		//err = me.EnsureNotNil()
		//if err != nil {
		//	break
		//}

		j, err = json.Marshal(me)
		if err != nil {
			break
		}
	}

	return MessageText(j)
}

func (me *MessageText) ToMessage() (*Message, error) {

	var err error
	var ret Message

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = json.Unmarshal(me.ByteArray(), &ret)
		if err != nil {
			break
		}

		err = ret.Validate()
	}

	return &ret, err
}

