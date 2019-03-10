package gearbox

import (
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

const ProjectFilename = "gearbox.json"
const ProjectFileHelpUrl = "https://docs.gearbox.works/projects/gearbox.json"

var _ util.FilepathHelpUrlGetter = (*ProjectFile)(nil)

type ProjectFile struct {
	JsonVersion JsonVersion `json:"gearbox"`
	Hostname    string      `json:"hostname"`
	Aliases     Aliases     `json:"aliases"`
	StackBag    StackBag    `json:"stack"`
	stackMap    StackMap    `json:"-"`
	Filepath    string      `json:"-"`
}

type projectDetails struct {
	Aliases  Aliases  `json:"aliases"`
	StackMap StackMap `json:"stack"`
	Filepath string   `json:"filepath"`
}

type JsonMeta struct {
	JsonVersion string   `json:"json"`
	Version     string   `json:"version"`
	Website     string   `json:"website"`
	Readme      []string `json:"readme"`
}

func (me *ProjectFile) ExportProjectDetails() *projectDetails {
	return &projectDetails{
		Aliases:  me.Aliases,
		StackMap: me.stackMap,
		Filepath: me.Filepath,
	}
}

func (me *ProjectFile) GetFilepath() string {
	return me.Filepath
}
func (me *ProjectFile) GetHelpUrl() string {
	return ProjectFileHelpUrl
}

func (me *ProjectFile) Unmarshal(j []byte) (status Status) {
	for range only.Once {
		err := util.UnmarshalJson(j, me)
		if err != nil {
			status = NewStatus(&StatusArgs{
				HelpfulError: err.(util.HelpfulError),
				HttpStatus:   http.StatusInternalServerError,
			})
			break
		}
		status = NewOkStatus("bytes unmarshalled")
	}
	return status
}

func (me *ProjectFile) FixupStackMap() (status Status) {
	//for role,item := range me.StackBag {
	//	var svc := &Service{}
	//	if sr,ok := item.(string); ok {
	//
	//	}
	//	err := util.UnmarshalJson(j, me)
	//	if err != nil {
	//		status = NewStatus(&StatusArgs{
	//			HelpfulError: err.(util.HelpfulError),
	//			HttpStatus:   http.StatusInternalServerError,
	//		})
	//		break
	//	}
	//	status = NewOkStatus("bytes unmarshalled")
	//}
	return status
}
