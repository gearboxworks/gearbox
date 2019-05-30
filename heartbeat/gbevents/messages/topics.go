package messages

import (
	"fmt"
	"gearbox/help"
	"gearbox/only"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"strings"
)

type Topic struct {
	Address  MessageAddress
	SubTopic
}
type Topics []Topic


////////////////////////////////////////////////////////////////////////////////
type TopicString string

func (me *TopicString) String() string {

	return string(*me)
}

func (me *TopicString) Split() Topic {

	var topic Topic

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

func StringToTopic(t string) Topic {

	ts := TopicString(t)

	return ts.Split()
}

func StringsToTopic(client string, topic string) Topic {

	ts := TopicString(fmt.Sprintf(TopicSeparator, client, topic))

	return ts.Split()
}

func CreateTopicGlob(client MessageAddress) *Topic {

	return &Topic{Address: client, SubTopic: "*"}
}

func CreateTopic(client MessageAddress, topic SubTopic) *Topic {

	return &Topic{Address: client, SubTopic: topic}
}

const TopicSeparator = "/%s/%s"
func SprintfTopic(address MessageAddress, topic SubTopic) string {

	return fmt.Sprintf(TopicSeparator, address.String(), topic.String())
}



////////////////////////////////////////////////////////////////////////////////
/*
func (me *Topic) CreateTopic(address MessageAddress) {

	foo := fmt.Sprintf("%s/%s", address, me.String())

	*me = Topic(foo)
}

func CreateTopic(id string) (string) {

	var te Topic
	te.CreateTopic(id)

	return te.String()

}
*/

func (me *Topic) EnsureNotNil() status.Status {

	var sts status.Status

	switch {
		case me == nil:
			sts = status.Warn("").
				SetMessage("topic is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		case me.Address.IsNil() == true:
			sts = status.Warn("").
				SetMessage("topic address is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		case me.SubTopic.EnsureNotNil().IsError():
			sts = me.SubTopic.EnsureNotNil()

		default:
			sts = status.Success("topic not nil")
	}

	return sts
}

func (me *Topic) String() string {

	var sts status.Status
	var s string

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		s = fmt.Sprintf(TopicSeparator, me.Address.String(), me.SubTopic.String())
	}

	return s
}


////////////////////////////////////////////////////////////////////////////////
type SubTopic string
var IsEmptySubTopic SubTopic
func (me *SubTopic) EnsureNotNil() status.Status {

	var sts status.Status

	switch {
		case me == nil:
			sts = status.Warn("").
				SetMessage("subtopic is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		default:
			sts = status.Success("subtopic not nil")
	}

	return sts
}

func (me *SubTopic) String() string {

	return string(*me)
}

// var AllTopics = Topic("*")
