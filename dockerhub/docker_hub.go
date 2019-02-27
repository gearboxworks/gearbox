package dockerhub

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
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

func (me *DockerHub) RequestAvailableContainerNames(query ...*ContainerQuery) (cns ContainerNames) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}

	for range only.Once {

		cns = make(ContainerNames, 0)

		crl, _, err := util.HttpGet(ContainerRespositoryListUrl)
		if err != nil {
			log.Printf("failed to download list of container repositories from '%s': %s ",
				ContainerRespositoryListUrl,
				err.Error(),
			)
			break
		}
		ar := ApiResponse{}
		err = json.Unmarshal(crl, &ar)
		if err != nil {
			log.Printf("failed to unmarshal response from '%s': %s ",
				ContainerRespositoryListUrl,
				err.Error(),
			)
			break
		}
		rb, err := json.Marshal(ar.Results)
		if err != nil {
			log.Printf("failed to marshal list of container repositories from '%s': %s ",
				ContainerRespositoryListUrl,
				err.Error(),
			)
			break
		}
		rs := Repositories{}
		err = json.Unmarshal(rb, &rs)
		if err != nil {
			log.Printf("failed to unmarshal list of container repositories from '%s': %s ",
				ContainerRespositoryListUrl,
				err.Error(),
			)
			break
		}
		for _, r := range rs {
			rc, err := me.RequestRepositoryContainers(ContainerName(r.Name), _query)
			if err != nil {
				continue
			}
			for _, c := range rc {
				cns = append(cns, ContainerName(fmt.Sprintf("%s/%s/%s", r.User, r.Name, c.Name)))
			}
		}
	}
	return cns
}

func (me *DockerHub) RequestRepositoryContainers(name ContainerName, query ...*ContainerQuery) (cs Containers, err error) {
	var _query *ContainerQuery
	if len(query) == 0 {
		_query = &ContainerQuery{}
	} else {
		_query = query[0]
	}
	for range only.Once {
		cs = make(Containers, 0)

		url := strings.Replace(AvailableContainerListUrl, "{name}", string(name), 1)

		acl, _, err := util.HttpGet(url)
		if err != nil {
			log.Printf("failed to download list of containers for '%s' from '%s': %s ",
				name,
				url,
				err.Error(),
			)
			break
		}
		ar := ApiResponse{}
		err = json.Unmarshal(acl, &ar)
		if err != nil {
			log.Printf("failed to unmarshal response from '%s': %s ",
				url,
				err.Error(),
			)
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
			log.Printf("failed to unmarshal list of containers from '%s': %s ",
				url,
				err.Error(),
			)
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
					break
				}
				jv, err := version.NewVersion(string(cs[j].Name))
				if err != nil {
					break
				}
				lt = iv.LessThan(jv)
			}
			return lt
		})
	}
	return cs, err
}
