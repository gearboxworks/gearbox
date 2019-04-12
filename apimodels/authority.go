package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

const AuthorityType = "authority"

var NilAuthority = (*Authority)(nil)
var _ apimodeler.ApiItemer = NilAuthority

type AuthorityMap map[types.Stackname]*Authority
type Authorities []*Authority

type Authority struct {
	AuthorityId types.AuthorityDomain `json:"authority_id"`
}

func NewFromGearsAuthority(ctx *apimodeler.Context, authority types.AuthorityDomain) (gs *Authority, sts status.Status) {
	return NewAuthority(authority), sts
}

func NewAuthority(authority types.AuthorityDomain) *Authority {
	return &Authority{
		AuthorityId: authority,
	}
}

func (me *Authority) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *Authority) GetType() apimodeler.ItemType {
	return AuthorityType
}

func (me *Authority) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *Authority) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.AuthorityId)
}

func (me *Authority) SetId(itemid apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		if !strings.Contains(string(itemid), ".") {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("authority domain '%s' does not contain a period ('.')", itemid),
			})
			break
		}
		me.AuthorityId = types.AuthorityDomain(itemid)
	}
	return sts
}

func (me *Authority) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}

func (me *Authority) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.ApiItemer) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}
