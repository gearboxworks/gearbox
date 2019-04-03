package docker

import (
	"encoding/json"
	"io/ioutil"
)

func GetRemoteImageTagTree(token, image string) (TagTree, error) {
	tags, err := GetRemoteImageTags(token, image)
	if err != nil {
		return TagTree{}, err
	}
	return NewTagTree(image, tags)
}

func GetRemoteImageTags(token, image string) (TagList, error) {

	au := BuildUrl(
		RegistryDomain,
		HubApiVersion,
		"/"+image+"/tags/list",
		EmptyQueryList,
	)

	client, req, err := NewHttpGetRequest(au)
	if err != nil {
		return TagList{}, err
	}
	resource.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return TagList{}, err
	}
	tags, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	tlr := TagListResponse{}
	err = json.Unmarshal(tags, &tlr)
	return tlr.Tags, err

}
