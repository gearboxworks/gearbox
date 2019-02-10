package gearbox

import (
	"fmt"
	"strings"
)

const ProjectFile = "project.json"

type Projects []*Project

type Project struct {
	Name    string
	Domain  string
	Enabled bool
}

func NewProject(name string) *Project {
	domain := name
	if !strings.Contains(name, ".") {
		domain = fmt.Sprintf("%s.local", name)
	}
	pr := Project{
		Name:   name,
		Domain: domain,
	}
	return &pr
}
