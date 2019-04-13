package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

const AuthorityModelType = "authority"

var NilAuthorityModel = (*AuthorityModel)(nil)
var _ apimodeler.Itemer = NilAuthorityModel

type AuthorityModelMap map[types.Stackname]*AuthorityModel
type AuthorityModels []*AuthorityModel

type AuthorityModel struct {
	AuthorityId types.AuthorityDomain `json:"authority_id"`
}

func NewFromGearsAuthority(ctx *apimodeler.Context, authority types.AuthorityDomain) (gs *AuthorityModel, sts status.Status) {
	return NewAuthority(authority), sts
}

func NewAuthority(authority types.AuthorityDomain) *AuthorityModel {
	return &AuthorityModel{
		AuthorityId: authority,
	}
}

func (me *AuthorityModel) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *AuthorityModel) GetType() apimodeler.ItemType {
	return AuthorityModelType
}

func (me *AuthorityModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *AuthorityModel) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.AuthorityId)
}

func (me *AuthorityModel) SetId(itemid apimodeler.ItemId) (sts status.Status) {
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

func (me *AuthorityModel) GetItem() (apimodeler.Itemer, status.Status) {
	return me, nil
}

func (me *AuthorityModel) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.Itemer) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}
