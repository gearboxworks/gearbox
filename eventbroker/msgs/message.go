package msgs

import (
	"encoding/json"
	"fmt"
	"github.com/gearboxworks/go-status/only"
)

type Message struct {
	Time   MessageTime
	Source Address
	Topic  Topic
	Text   Text

	//PayLoad
}

type PayLoad struct {
	Topic Topic
	Text  Text
}

func (me *Message) String() string {

	return fmt.Sprintf(`Time:%d  Source:%s  Topic:%s  Text:%s`,
		me.Time.Unix(),
		me.Source.String(),
		me.Topic.String(),
		me.Text.String(),
	)
}

func (me *Message) Validate() error {

	var err error

	for range only.Once {
		//err = me.Text.EnsureNotEmpty()
		//if err != nil {
		//	break
		//}

		err = me.Topic.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.Source.EnsureNotEmpty()
		if err != nil {
			break
		}
	}

	return err
}

func (me *Message) ToMessageText() Text {

	var err error
	var j []byte

	for range only.Once {
		//err = me.EnsureNotEmpty()
		//if err != nil {
		//	break
		//}

		j, err = json.Marshal(me)
		if err != nil {
			break
		}
	}

	return Text(j)
}
