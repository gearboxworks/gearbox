package projectcfg

import (
	"encoding/json"
	"fmt"
	"gearbox/projectcfg/util"
	"gearbox/projectcfg/util/types"
	"io/ioutil"
	"log"
	"strings"
)

const DefaultFilename string = "project.json"

var Instance *ProjectCfg

func Initialize(filepath string) {
	Instance = Load(filepath)
}

type projectCfg struct {
	Version string `json:"version"`
}

type ProjectCfg struct {
	ProjectCfg  projectCfg             `json:"projectcfg"`
	Name        string                 `json:"name,omitempty"`
	Team        Team                   `json:"team,omitempty"`
	Description string                 `json:"description,omitempty"`
	Namespace   string                 `json:"namespace,omitempty"`
	Slug        string                 `json:"slug,omitempty"`
	Prefix      string                 `json:"prefix,omitempty"`
	Type        types.Project          `json:"type,omitempty"`
	Aliases     []string               `json:"aliases,omitempty"`
	Root        *string                `json:"-"`
	Dev         *Dev                   `json:"dev,omitempty"`
	Source      *Source                `json:"source,omitempty"`
	Deploy      *Deploy                `json:"deploy,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
	//Stack            Stack                  `json:"stack,omitempty"`
	//Hosts            HostMap                `json:"hosts,omitempty"`
	//loadedVendors VendorMap
}

func Load(cfgPath string) *ProjectCfg {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		util.Error(err)
	}
	p := ProjectCfg{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		log.Fatal(err)
	}
	return &p
}

func New(name string, root *string) *ProjectCfg {
	domain := name
	if !strings.Contains(name, ".") {
		domain = fmt.Sprintf("%s.local", name)
	}
	pr := ProjectCfg{
		//		Root:     root,
		Name:     name,
		Hostname: domain,
	}
	return &pr
}

func (me *ProjectCfg) String() string {
	j, err := json.MarshalIndent(me, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(j)
}

//func (me *ProjectCfg) GetVendor(name VendorName) (fw Vendor) {
//	fw, _ = me.loadedVendors[name]
//	return
//}

//
//@TODO Needs to allow additonal vendor sources defined in the config file
//      Sources need to be able to be downloaded, but need to be secure
//@TODO Waiting on a GoLang bug fix: https://youtrack.jetbrains.com/issue/GO-6289
//
//var pathTemplate = "%s/vendors/*.pcfw"
//
//func (me *ProjectCfg) LoadVendors() {
//	path := fmt.Sprintf(pathTemplate,util.GetProjectDir())
//	files, err := filepath.Glob(path)
//	if err != nil {
//		log.Fatal(fmt.Sprintf("Directory listing of vendor file in '%s' failed: %s",path,err))
//	}
//	for _,file := range files {
//		fwpi, err := plugin.Open(file)
//		if err != nil {
//			log.Fatal(fmt.Sprintf("Open vendor file '%s' failed: %s",file,err))
//		}
//		loader, err := fwpi.Lookup("GetInstance")
//		if err != nil {
//			log.Fatal(fmt.Sprintf("Lookup GetInstance in vendor file '%s' failed: %s",file,err))
//		}
//		getter, ok := loader.(func() interface{})
//		if !ok {
//			log.Fatal(fmt.Sprintf("Type assert for vendor file '%s' failed: %s",file,err))
//		}
//		generic := getter()
//		vendor, ok := generic.(Vendor)
//		if !ok {
//			log.Fatal(fmt.Sprintf("Vendor file '%s' in not a valid vendor: %s",file,err))
//		}
//		me.loadedVendors[VendorName(filepath.Base(file))] = vendor
//	}
//}
