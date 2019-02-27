package util

import (
	"gearbox/only"
	"io/ioutil"
	"log"
	"net/http"
)

func CloseResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Printf(
			"Could not close response body from HttpGet: %s\n",
			err.Error(),
		)
	}
}

func HttpGet(url string) (body []byte, statusCode int, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("# Http Get - Panic occurred for '%s': %s\n", url, r)
		}
	}()

	for range only.Once {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("# Error in HttpGet: '%s'\n", err)
			break
		}
		//request.SetBasicAuth(_api.User, _api.Password)
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		cli := &http.Client{}
		response, err := cli.Do(request)
		statusCode = response.StatusCode
		if err != nil {
			log.Printf("# Error in HttpGet: '%s'\n", err)
			break
		}
		defer CloseResponseBody(response)
		body, err = ioutil.ReadAll(response.Body)
		switch statusCode {
		case 200:
			// OK
			// bodyBytes, err2 := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("# Api Get - Nothing returned for url: %s\n", url)
				break
			}

		case 401:
			// Invalid login.
			log.Fatal("# Api Get - Requires correct API keys.\n")
			break

		case 403:
			// Permission denied?
			log.Printf("# Api Get - Permission denied for url: %s\n", url)
			log.Printf("# Api Get - Return page: '%s'\n", body)
			break

		default:
			// 404 error
			log.Printf("# Api Get - HTML %d returned from url:%s.\n", response.StatusCode, url)
			log.Printf("# Api Get - Return page: '%s'\n", body)
			break
		}

	}
	return body, statusCode, err
}
