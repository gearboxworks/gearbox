package test

import (
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient struct {
	response *http.Response
	body     []byte
	url      string
	status   stat.Status
}

func (me *HttpClient) GET(url string) (status stat.Status) {
	me.url = url
	r, err := http.Get(url)
	if err != nil {
		status = stat.NewErrorStatus(err, &stat.Args{
			Message: fmt.Sprintf("error while requesting '%s'", url),
		})
	}
	me.response = r
	me.status = status
	return status
}

func (me *HttpClient) GetBody() (body []byte, status stat.Status) {
	for range only.Once {
		if me.response == nil {
			log.Fatal("cannot call HttpClient.GetBody() before response is set")
		}
		var err error
		body, err = ioutil.ReadAll(me.response.Body)
		if err != nil {
			status = stat.NewErrorStatus(err, &stat.Args{
				Message: fmt.Sprintf("error reading response body for '%s'", me.url),
				Data:    me.response,
			})
		}
	}
	me.body = body
	return body, status
}
