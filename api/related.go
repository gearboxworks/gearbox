package api

type RelatedMap map[RouteName]*Related

type Related struct {
	List   RouteName
	New    RouteName
	Update RouteName
	Delete RouteName
	Item   RouteName
	Others RouteNameMap
}

func (me *Api) Relate(primary RouteName, related *Related) {

	rtm, exists := me.RelatedMap[primary]
	if !exists {
		me.RelatedMap[primary] = make(RouteNameMap, 0)
		rtm = me.RelatedMap[primary]
	}

	switch primary {

	case related.List:
		if related.Item != "" {
			rtm[related.Item] = ItemResource
		}
		if related.New != "" {
			rtm[related.New] = MakeNewResource
		}

	case related.Item:
		if related.List != "" {
			rtm[related.List] = ListResource
			me.Relate(related.List, &Related{
				Item:   primary,
				List:   related.List,
				New:    related.New,
				Others: related.Others,
			})
		}
		if related.New != "" {
			rtm[related.New] = MakeNewResource
		}
		if related.Update != "" {
			rtm[related.Update] = UpdateResource
		}
		if related.Delete != "" {
			rtm[related.Delete] = DeleteResource
		}

	}
	if related.Others != nil {
		for o, t := range related.Others {
			rtm[o] = t
		}
	}
	me.RelatedMap[primary] = rtm
}
