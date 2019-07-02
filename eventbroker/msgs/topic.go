package msgs

import (
	"errors"
	"fmt"
)

type Topics []Topic
type Topic struct {
	Address Address
	SubTopic
}

func NewGlobTopic(a Address) *Topic {
	return NewTopic(a, "*")
}

func NewTopic(a Address, st SubTopic) *Topic {
	return &Topic{
		Address:  a,
		SubTopic: st,
	}
}

func (me *Topic) EnsureNotNil() (err error) {

	switch {
	case me == nil:
		err = errors.New("topic is nil")

	case me.Address.EnsureNotEmpty() != nil:
		err = errors.New("topic address is empty")

	case me.SubTopic.EnsureNotEmpty() != nil:
		err = me.SubTopic.EnsureNotEmpty()

	}

	return err
}

func (me *Topic) String() string {

	var err error
	var s string

	err = me.EnsureNotNil()
	if err != nil {
		return s
	}

	return fmt.Sprintf(TopicPattern,
		me.Address.String(),
		me.SubTopic.String())
}
