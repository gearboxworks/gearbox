package msgs

import "time"

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
