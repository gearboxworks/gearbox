package dockerhub

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"gearbox/util"
	"github.com/hashicorp/go-version"
	"log"
	"sort"
	"strings"
)

type ExternalUrlId string

const (
	ContainerRespositoryListUrl = "https://hub.docker.com/v2/repositories/wplib/?page_size=256"
	AvailableContainerListUrl   = "https://hub.docker.com/v2/repositories/wplib/{name}/tags/?page_size=256"
)

type DockerHub struct {
}

type ContainerQuery struct {
	IncludePatches bool
}

func (me *DockerHub) RequestAvailableContainerNames(query ...*ContainerQuery) (cns ContainerNames, status stat.Status) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}

	for range only.Once {

		cns = make(ContainerNames, 0)

		crl, _, status := util.HttpGet(ContainerRespositoryListUrl)
		if status.IsError() {
			status.PriorStatus = status.String()
			status.Message = fmt.Sprintf("failed to download list of container repositories from '%s'",
				ContainerRespositoryListUrl,
			)
			status.Help = "check your Internet access, or " + stat.ContactSupportHelp()
			break
		}
		ar := ApiResponse{}
		err := json.Unmarshal(crl, &ar)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("unable to unmarshal response from '%s'",
					ContainerRespositoryListUrl,
				),
			})
			break
		}
		rb, err := json.Marshal(ar.Results)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("unable to marshal container repositories from '%s'",
					ContainerRespositoryListUrl,
				),
			})
			break
		}
		rs := Repositories{}
		err = json.Unmarshal(rb, &rs)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("unable to unmarshal container repositories from '%s'",
					ContainerRespositoryListUrl,
				),
			})
			break
		}
		for _, r := range rs {
			rc, status := me.RequestRepositoryContainers(ContainerName(r.Name), _query)
			if status.IsError() {
				// @TODO Should we do anything here?!?
				continue
			}
			for _, c := range rc {
				cns = append(cns, ContainerName(fmt.Sprintf("%s/%s/%s", r.User, r.Name, c.Name)))
			}
		}
	}
	return cns, status
}

func (me *DockerHub) RequestRepositoryContainers(name ContainerName, query ...*ContainerQuery) (cs Containers, status stat.Status) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}
	for range only.Once {
		cs = make(Containers, 0)

		url := strings.Replace(AvailableContainerListUrl, "{name}", string(name), 1)

		acl, _, status := util.HttpGet(url)
		if status.IsError() {
			status.PriorStatus = status.String()
			status.Message = fmt.Sprintf("failed to download list of containers for '%s' from '%s'",
				name,
				url,
			)
			break
		}
		ar := ApiResponse{}
		err := json.Unmarshal(acl, &ar)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("unable to unmarshal response from '%s'",
					ContainerRespositoryListUrl,
				),
			})
			break
		}
		cb, err := json.Marshal(ar.Results)
		if err != nil {
			log.Printf("failed to marshal list of containers from '%s': %s ",
				ContainerRespositoryListUrl,
				err.Error(),
			)
			break
		}

		err = json.Unmarshal(cb, &cs)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("unable to marshal container repositories from '%s'",
					ContainerRespositoryListUrl,
				),
			})
			break
		}

		if !_query.IncludePatches {
			for i := len(cs) - 1; i >= 0; i-- {
				if strings.Count(string(cs[i].Name), ".") > 1 {
					cs = RemoveContainer(i, cs)
				}
			}
		}

		for i := len(cs) - 1; i >= 0; i-- {
			if cs[i].Name == "latest" {
				cs = RemoveContainer(i, cs)
				break
			}
		}

		sort.Slice(cs, func(i, j int) (lt bool) {
			for range only.Once {
				iv, err := version.NewVersion(string(cs[i].Name))
				if err != nil {
					status = stat.NewFailedStatus(&stat.Args{
						Error: err,
					})
					break
				}
				jv, err := version.NewVersion(string(cs[j].Name))
				if err != nil {
					status = stat.NewFailedStatus(&stat.Args{
						Error: err,
					})
					break
				}
				lt = iv.LessThan(jv)
			}
			return lt
		})
	}
	return cs, status
}
