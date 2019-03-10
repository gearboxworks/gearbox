package cardinality

import (
	"fmt"
	"github.com/projectcfg/projectcfg/util/only"
	"strings"
)

type Type string

const (
	ZeroOrOne  Type = "zero-or-one"
	ExactlyOne Type = "exactly-one"
	One        Type = "one"
	Many       Type = "many"
	//ToOne      Type = "to-one"
	//ToMany     Type = "to-many"
)

type Cardinality struct {
	From Type `json:"from"`
	To   Type `json:"to"`
}

func NewCardinality(from, to Type) (c *Cardinality, err error) {
	c = &Cardinality{
		From: from,
		To:   to,
	}
	_, err = c.String()
	return c, err
}

func (me *Cardinality) String() (cs string, err error) {
	var ok bool
	for range only.Once {
		cs = fmt.Sprintf("%s-to-%s", me.From, me.To)
		cs = strings.Replace("to-to", "to", cs, 1)
		if cs[0] == 't' {
			break
		}
		if strings.HasSuffix(cs, string(ZeroOrOne)) {
			break
		}
		if strings.HasSuffix(cs, string(ExactlyOne)) {
			break
		}
		ok = true
	}
	if !ok {
		err = fmt.Errorf("invalid cardinality '%s'", cs)
	}
	return cs, err
}
