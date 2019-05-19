package config

import "gearbox/types"

type NicknameMap map[types.AbsoluteDir]types.Nickname

func (me *Config) GetNicknameMap() (nnm NicknameMap) {
	bdm := me.GetBasedirMap()
	nnm = make(NicknameMap, len(bdm))
	for nn, bd := range bdm {
		nnm[bd.Basedir] = nn
	}
	return nnm
}
