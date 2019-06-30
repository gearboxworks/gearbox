package msgs

import (
	"errors"
)

type Addresses []Address
type Address string

func (me Address) MakeMessage(to Address, subtopic SubTopic, text Text) Message {
	return Message{
		Source: me,
		Text:   text,
		Topic:  me.MakeTopic(to, subtopic),
	}
}

func (me Address) MakeTopic(to Address, subtopic SubTopic) Topic {
	return Topic{
		Address:  to,
		SubTopic: subtopic,
	}
}

func (me Address) String() string {
	return string(me)
}

func (me Address) EnsureNotEmpty() (err error) {
	if me == "" {
		err = errors.New("message address is empty")
	}
	return err
}
