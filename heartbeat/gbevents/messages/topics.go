package messages

import (
	"fmt"
	"gearbox/global"
)


type Topic string

func (me *Topic) CreateTopic(id string) {

	foo := fmt.Sprintf("%s/%s", global.Brandname, id)

	*me = Topic(foo)
}

func CreateTopic(id string) (string) {

	var te Topic
	te.CreateTopic(id)

	return te.ToString()

}

func (me *Topic) ToString() (string) {

	return string(*me)
}

