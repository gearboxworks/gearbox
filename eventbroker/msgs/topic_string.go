package msgs

import "strings"

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
	case ar[0] != "":
		fallthrough
	case len(ar) <= 2:
		// Failed

	case (ar[0] == "") && (len(ar) > 2):
		// If first element is "", then we have started with a '/'.
		topic.Address = Address(ar[1])
		topic.SubTopic = SubTopic(strings.Join(ar[2:], "/"))
	}

	return topic
}

func StringToTopic(t string) Topic {

	ts := TopicString(t)

	return ts.Split()
}
