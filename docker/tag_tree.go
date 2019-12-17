package docker

import (
	"sort"
	"strconv"
	"strings"
)

type TagTreeTraverser interface {
	Tag() string
	Children() *map[string]TagTree
}

type TagTree struct {
	tag      string
	children map[string]TagTree
	sorted   []TagTree
}

func (tt TagTree) Tag() string {
	return tt.tag
}

func (tt TagTree) Children() *map[string]TagTree {
	return &tt.children
}

func (tt TagTree) SortedChildren() []TagTree {
	return tt.sorted
}

func (tt TagTree) GetChild(tag string) *TagTree {
	tc := tt.children[tag]
	return &tc
}

func NewEmptyTagTree(tag string) *TagTree {
	return &TagTree{
		tag:      tag,
		children: make(map[string]TagTree),
	}
}

func NewTagTree(name string, list TagList) (TagTree, error) {
	tt := NewEmptyTagTree(name)
	for _, t := range list {
		if t == "latest" {
			continue
		}
		ttt := tt
		nums := strings.Split(t, ".")
		var ok bool
		for _, n := range nums {
			c := *ttt.Children()
			if _, ok = c[n]; ! ok {
				ttt.children[n] = *NewEmptyTagTree(n)
			}
			ttt = ttt.GetChild(n)
		}
	}
	return tt.getSorted(), nil
}

func (t TagTree) getSorted() TagTree {
	return *&TagTree{
		tag:t.tag,
		children:t.children,
		sorted:t.getSortedChildren(),
	}
}

func (t TagTree) getSortedChildren() []TagTree {
	sorted := *new([]TagTree)
	pl := SortablePairList{}
	i := 0
	for _, tag := range *t.Children() {
		pl = append(pl, SortablePair{
			Key:    i,
			Value:  tag,
			Cargo:  &t,
			IsLess: isLess,
		})
		i++
	}
	sort.Sort(pl)
	for _, kv := range pl {
		tt:= kv.Value.(TagTree)
		tt.sorted = tt.getSortedChildren()
		sorted = append(sorted,tt)
	}
	return sorted
}

func isLess(i, j interface{}) bool {

	isp := i.(SortablePair)
	iv, err := strconv.Atoi(isp.Value.(TagTree).Tag())
	if err != nil {
		return false
	}
	jsp := j.(SortablePair)
	jv, err := strconv.Atoi(jsp.Value.(TagTree).Tag())
	if err != nil {
		return false
	}

	return iv < jv;
}

/*
dd is variadic but only so I can make it optional.
The first one passed is the only one used.
 */
func (t TagTree) Print(dd ...byte) {
	var d int
	if len(dd) == 0 {
		d = 0
	} else {
		d = int(dd[0])
	}
	println(t.Tag())
	for _, ct := range t.sorted {
		print(strings.Repeat("-", d*3))
		ct.Print(byte(d + 1))
	}

}
