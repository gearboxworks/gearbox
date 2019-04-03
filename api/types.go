package api

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"strings"
)

type ResourceType string

type TemplateVars []TemplateVar

func (me TemplateVars) Values() (values []string) {
	values = make([]string, len(me))
	for i, s := range me {
		values[i] = string(s)
	}
	return values
}

type TemplateVar string

// TemplateVar and ResourceVarName are basically the same
// thing but discovered in two different contexts.  Might
// merge them into one later.
type ResourceVarName string

type ResourcesMap map[RouteName]RouteNameMap
type RouteNameMap map[RouteName]ResourceType
type RouteName string

//type Links map[RouteName]UriTemplate
type Links map[RouteName]interface{}

type ValuesFuncMap map[RouteName]ValuesFunc
type ValuesFunc func(...interface{}) (ValuesFuncValues, status.Status)

type ValuesFuncValues []ValueFuncVarsValues
type ValueFuncVarsValues []string

type UriTemplates []UriTemplate
type UriTemplate string

func (me UriTemplate) Expand(vars UriTemplateVars) UriTemplate {
	url := me
	for _, v := range vars {
		url = UriTemplate(strings.Replace(string(url), fmt.Sprintf("{%s}", v.Name), v.Value, -1))
	}
	return url
}

func (me UriTemplate) Convert() (url UriTemplate) {
	for range only.Once {
		url = me
		if !strings.Contains(string(me), ":") {
			break
		}
		parts := strings.Split(string(me), "/")
		for i, p := range parts {
			if len(p) == 0 {
				continue
			}
			if []byte(p)[0] != ':' {
				continue
			}
			parts[i] = fmt.Sprintf("{%s}", p[1:])
		}
		url = UriTemplate(strings.Join(parts, "/"))
	}
	return url
}

type UriTemplateVars []*UriTemplateVar
type UriTemplateVar struct {
	Name  ResourceVarName
	Value string
}

func (me UriTemplateVars) Values() (svars []string) {
	svars = make([]string, len(me))
	i := 0
	for _, v := range me {
		svars[i] = v.Value
		i++
	}
	return svars
}

func (me UriTemplateVars) String() string {
	a := make([]string, len(me))
	i := 0
	for _, tv := range me {
		a[i] = tv.Value
		i++
	}
	return fmt.Sprintf("[%s]", strings.Join(a, "]["))
}

type EndpointMap map[RouteName]*Endpoint
type Endpoint struct {
	UriTemplate UriTemplate `json:"uri_template"`
	Methods     Methods     `json:"methods"`
}

type ResourceMap map[RouteName]UriTemplate
type MethodMap map[Method]ResourceMap
type Methods []Method
type Method string

type HandlerFunc func(*RequestContext) error

type UpstreamHandlerFunc func(*RequestContext) interface{}
