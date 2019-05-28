package messages

import "time"

type Message struct {
	Src string
	Time time.Time
	Text string
}
