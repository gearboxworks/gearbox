package gearbox

type StackMemberName string

type StackMemberMap map[StackMemberName]*StackMember

type StackMembers []*StackMember

type StackMember struct {
	Name       StackMemberName `json:"name"`
	Label      string          `json:"label"`
	ShortLabel string          `json:"short_label"`
	Examples   []string        `json:"examples"`
	StackName  StackName       `json:"stack"`
	MemberType string          `json:"member_type"`
	Optional   bool            `json:"optional"`
}
type StackMemberArgs StackMember

func (me *StackMember) String() string {
	return me.Label
}
