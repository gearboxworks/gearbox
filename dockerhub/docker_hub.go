package dockerhub

import (
	"encoding/json"
	"fmt"
	"gearbox/help"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/only"
	"github.com/hashicorp/go-version"
	"log"
	"sort"
	"strings"
)

type ExternalUrlId string

const (
	ContainerRepositoryListUrl = "https://hub.docker.com/v2/repositories/wplib/?page_size=256"
	AvailableContainerListUrl  = "https://hub.docker.com/v2/repositories/wplib/{name}/tags/?page_size=256"
)

type DockerHub struct {
}

type ContainerQuery struct {
	IncludePatches bool
}

func (me *DockerHub) RequestAvailableContainerNames(query ...*ContainerQuery) (cns ContainerNames, sts status.Status) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}

	for range only.Once {

		cns = make(ContainerNames, 0)

		crl, _, sts := util.HttpRequest(ContainerRepositoryListUrl)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message: fmt.Sprintf("failed to download list of container repositories from '%s'",
					ContainerRepositoryListUrl,
				),
				Help: "check your Internet access, or " + help.ContactSupportHelp(),
			})
			break
		}
		ar := ApiResponse{}
		err := json.Unmarshal(crl, &ar)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal response from '%s'",
					ContainerRepositoryListUrl,
				),
			})
			break
		}
		rb, err := json.Marshal(ar.Results)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marshal container repositories from '%s'",
					ContainerRepositoryListUrl,
				),
			})
			break
		}
		rs := Repositories{}
		err = json.Unmarshal(rb, &rs)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal container repositories from '%s'",
					ContainerRepositoryListUrl,
				),
			})
			break
		}
		for _, r := range rs {
			rc, sts := me.RequestRepositoryContainers(ContainerName(r.Name), _query)
			if status.IsError(sts) {
				// @TODO Should we do anything here?!?
				continue
			}
			for _, c := range rc {
				cns = append(cns, ContainerName(fmt.Sprintf("%s/%s/%s", r.User, r.Name, c.Name)))
			}
		}
	}
	return cns, sts
}

func (me *DockerHub) RequestRepositoryContainers(name ContainerName, query ...*ContainerQuery) (cs Containers, sts status.Status) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}
	for range only.Once {
		cs = make(Containers, 0)

		url := strings.Replace(AvailableContainerListUrl, "{name}", string(name), 1)

		acl, _, sts := util.HttpRequest(url)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message: fmt.Sprintf("failed to download list of containers for '%s' from '%s'",
					name,
					url,
				),
			})
			break
		}
		ar := ApiResponse{}
		err := json.Unmarshal(acl, &ar)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to unmarshal response from '%s'",
					ContainerRepositoryListUrl,
				),
			})
			break
		}
		cb, err := json.Marshal(ar.Results)
		if err != nil {
			log.Printf("failed to marshal list of containers from '%s': %s ",
				ContainerRepositoryListUrl,
				err.Error(),
			)
			break
		}

		err = json.Unmarshal(cb, &cs)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("unable to marshal container repositories from '%s'",
					ContainerRepositoryListUrl,
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
					sts = status.Wrap(err)
					break
				}
				jv, err := version.NewVersion(string(cs[j].Name))
				if err != nil {
					sts = status.Wrap(err)
					break
				}
				lt = iv.LessThan(jv)
			}
			return lt
		})
	}
	return cs, sts
}
