package types

type Nicknames []Nickname
type Nickname string

func (me Nicknames) Strings() []string {
	bdnns := make([]string, len(me))
	for i, bdnn := range me {
		bdnns[i] = string(bdnn)
	}
	return bdnns
}
