package docker

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
)

type AuthToken struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IssuedAt    string `json:"issued_at"`
}

func GetAuthToken(repo, scope string) (string, error) {
	p := "/token"
	q := url.Values{
		"service": {"registry.docker.io"},
		"scope":   {"repository:" + repo + ":" + scope},
	}
	au := BuildUrl(AuthDomain, "", p, q.Encode())
	client, req, err := NewHttpGetRequest(au)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	authToken := AuthToken{}
	err = json.Unmarshal(body, &authToken)
	return authToken.Token, err
}
