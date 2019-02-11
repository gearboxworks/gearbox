package gearbox

import (
	"encoding/json"
	"fmt"
	"github.com/projectcfg/projectcfg"
	"github.com/zserge/webview"
	"log"
	"strings"
)

type Bridge struct {
	Webview webview.WebView
	Gearbox *Gearbox
}

type BridgeProject struct {
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	IsEnabled bool   `json:"isEnabled"`
}

func NewBrideProject(p *projectcfg.Project) *BridgeProject {
	return &BridgeProject{
		Name:      p.Name,
		Domain:    p.Hostname,
		IsEnabled: p.IsEnabled,
	}
}

func (me *Bridge) LoadProjects() {
	projects := me.Gearbox.Config.Projects
	ps := make([]string, len(projects))
	for i, p := range me.Gearbox.Config.Projects {
		pj, err := json.Marshal(NewBrideProject(p))
		if err != nil {
			log.Fatal(err)
		}
		ps[i] = string(pj)
	}
	js := fmt.Sprintf("gearbox.projects = [%s];", strings.Join(ps, ","))
	err := me.Webview.Eval(js)
	if err != nil {
		log.Fatal(err)
	}
}
