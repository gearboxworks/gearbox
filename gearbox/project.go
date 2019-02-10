package gearbox

import (
	"fmt"
	"strings"
)

const ProjectFile = "project.json"

type ProjectMap map[string]*Project
type Projects []*Project

type Project struct {
	Root      *string
	Name      string
	Domain    string
	IsEnabled bool
}

func NewProject(name string, root *string) *Project {
	domain := name
	if !strings.Contains(name, ".") {
		domain = fmt.Sprintf("%s.local", name)
	}
	pr := Project{
		Root:   root,
		Name:   name,
		Domain: domain,
	}
	return &pr
}
