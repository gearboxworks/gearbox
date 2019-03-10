// See https://www.youtube.com/watch?v=YRWCa84pykM

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"sort"
	//"github.com/davecgh/go-spew/spew"
)

type visitor struct {
	pkgDecl map[*ast.GenDecl]bool
	locals  map[string]int
	globals map[string]int
}

func newVisitor(f *ast.File) visitor {
	decls := make(map[*ast.GenDecl]bool)
	for _, decl := range f.Decls {
		v, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		decls[v] = true
	}
	return visitor{
		decls,
		make(map[string]int),
		make(map[string]int),
	}
}
func main() {
	fs := token.NewFileSet()
	locals := make(map[string]int)
	globals := make(map[string]int)
	for _, arg := range os.Args[1:] {
		f, err := parser.ParseFile(fs, arg, nil, parser.AllErrors)
		if err != nil {
			log.Printf("Could not parse %s: %v", arg, err)
		}
		v := newVisitor(f)
		ast.Walk(v, f)
		for k, v := range v.locals {
			locals[k] += v
		}
		for k, v := range v.globals {
			globals[k] += v
		}
		//spew.Dump(f)
	}
	log.Printf("Most Common Locals:")
	printTopFive(locals)
	log.Printf("Most Common Globals:")
	printTopFive(globals)
}

func printTopFive(counts map[string]int) {
	type pair struct {
		s string
		n int
	}
	pairs := make([]pair, 0, len(counts))
	for s, n := range counts {
		pairs = append(pairs, pair{s, n})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].n > pairs[j].n })
	for i := 0; i < len(pairs) && i < 5; i++ {
		fmt.Printf("%6d %s\n", pairs[i].n, pairs[i].s)
	}
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.AssignStmt:
		if d.Tok != token.DEFINE {
			return v
		}
		for _, name := range d.Lhs {
			countLocalIdent(v, name)
		}
	case *ast.RangeStmt:
		//fmt.Printf("%v %v\n",d.Key,d.Value)
		countLocalIdent(v, d.Key)
		countLocalIdent(v, d.Value)
	case *ast.GenDecl:
		if d.Tok != token.VAR {
			return v
		}
		for _, spec := range d.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if ! ok {
				continue
			}
			for _, name := range value.Names {
				if "_" == name.Name {
					break
				}
				if v.pkgDecl[d] {
					v.globals[name.Name]++
				} else {
					v.locals[name.Name]++
				}
			}
		}
	default:
		//spew.Dump(d)
	}
	return v
}

func countLocalIdent(v visitor, n ast.Node) {
	for {
		ident, ok := n.(*ast.Ident);
		if ! ok {
			break
		}
		if "_" == ident.Name {
			break
		}
		if ident.Obj == nil {
			break
		}
		//fmt.Printf( "ident.Obj.Pos: %d\n", ident.Obj.Pos() )
		//fmt.Printf( "ident.Pos: %d\n", ident.Pos() )
		if ident.Obj.Pos() != ident.Pos() {
			break
		}
		v.locals[ident.Name]++
		//fmt.Printf("lhs: %T\n",name)
		break
	}

}
