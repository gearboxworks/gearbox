package msgs

import "errors"

type SubTopics []SubTopic
type SubTopic string

func (me SubTopic) EnsureNotEmpty() (err error) {
	if me == "" {
		err = errors.New("subtopic is empty")
	}
	return err
}

func (me SubTopic) String() string {
	return string(me)
}
