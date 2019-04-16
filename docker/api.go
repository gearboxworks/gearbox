package docker

import (
	"gearbox/util"
	"net/http"
)

type HttpHeaders map[string]string

func NewHttpGetRequest(url string) (*http.Client, *http.Request, error) {
	return NewHttpRequest("GET", url)
}

func NewHttpRequest(method, url string) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err == nil {
		req.Header.Add("User-Agent", UserAgent)
		req.Header.Add("Content-Type", "application/json")
	}
	return client, req, err
}

func BuildUrl(host, version, path, query string) string {
	scheme := "https"
	if version != "" {
		version = "/v" + version
	}
	if util.CharAt(path, 0) != '/' {
		path = "/" + path
	}
	if query != "" {
		query = "?" + query
	}
	return scheme + "://" + host + version + path + query
}
