package gearbox

import (
	"bytes"
	"fmt"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/printer"
	"io/ioutil"
	"os"
)

type Php struct {
}

func (me *Php) Parse(filepath string) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Could not load %s", filepath)
		return
	}

	s := bytes.NewBufferString(string(buf))
	p := php7.NewParser(s, filepath)
	p.Parse()

	for _, e := range p.GetErrors() {
		fmt.Println(e)
	}

	rn := p.GetRootNode()

	file := os.Stdout
	pr := printer.NewPrinter(file)

	sl := rn.(*stmt.StmtList)

	pr.Print(sl)

}
