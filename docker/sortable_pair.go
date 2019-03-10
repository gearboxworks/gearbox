package docker

type SortablePair struct {
	Key int
	Value interface{}
	Cargo interface{}
	IsLess func(i,j interface{}) bool
}

type SortablePairList []SortablePair

func (p SortablePairList) Len() int {
	return len(p)
}

func (p SortablePairList) Swap(i, j int){
	p[i], p[j] = p[j], p[i]
}

func (p SortablePairList) Less(i, j int) bool {
	fn := p[i].IsLess
	return fn(p[i],p[j])

}
