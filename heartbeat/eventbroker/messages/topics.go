package messages

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/states"
	"strings"
)

type MessageTopic struct {
	Address  MessageAddress
	SubTopic
}
type Topics []MessageTopic


////////////////////////////////////////////////////////////////////////////////
type TopicString string

func (me *TopicString) String() string {

	return string(*me)
}

func (me *TopicString) Split() MessageTopic {

	var topic MessageTopic

	// Expecting a Topic like: "/address/topic/..."
	// This means that [0] == "", [1] == address, [2:] == true topic.
	ar := strings.Split(me.String(), "/")
	switch {
		case (ar[0] != ""):
			fallthrough
		case len(ar) <= 2:
			// Failed

		case (ar[0] == "") && (len(ar) > 2):
			// If first element is "", then we have started with a '/'.
			topic.Address = MessageAddress(ar[1])
			topic.SubTopic = SubTopic(strings.Join(ar[2:], "/"))
	}

	return topic
}

func StringToTopic(t string) MessageTopic {

	ts := TopicString(t)

	return ts.Split()
}

func StringsToTopic(client string, topic string) MessageTopic {

	ts := TopicString(fmt.Sprintf(TopicSeparator, client, topic))

	return ts.Split()
}


func CreateTopicGlob(client MessageAddress) *MessageTopic {

	return &MessageTopic{Address: client, SubTopic: states.ActionGlob}
}
func (client MessageAddress) CreateTopicGlob() *MessageTopic {

	return &MessageTopic{Address: client, SubTopic: states.ActionGlob}
}


func CreateTopic(client MessageAddress, topic SubTopic) *MessageTopic {

	return &MessageTopic{Address: client, SubTopic: topic}
}
func (client MessageAddress) CreateTopic(topic SubTopic) *MessageTopic {

	return &MessageTopic{Address: client, SubTopic: topic}
}


func SprintfTopic(address MessageAddress, topic SubTopic) string {

	return fmt.Sprintf(TopicSeparator, address.String(), topic.String())
}



////////////////////////////////////////////////////////////////////////////////

func (me *MessageTopic) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("topic is nil")

		case me.Address.EnsureNotNil() != nil:
			err = errors.New("topic address is nil")

		case me.SubTopic.EnsureNotNil() != nil:
			err = me.SubTopic.EnsureNotNil()

		default:
			err = nil
	}

	return err
}

func (me *MessageTopic) String() string {

	var err error
	var s string

	err = me.EnsureNotNil()
	if err != nil {
		return s
	}

	s = fmt.Sprintf(TopicSeparator, me.Address.String(), me.SubTopic.String())

	return s
}


////////////////////////////////////////////////////////////////////////////////
type SubTopic string
type SubTopics []SubTopic
func (me *SubTopic) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("subtopic is nil")

		default:
			err = nil
	}

	return err
}

func (me *SubTopic) String() string {

	return string(*me)
}

// var AllTopics = Topic("*")
