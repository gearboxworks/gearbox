package util

import (
	"fmt"
	"gearbox/only"
	"gearbox/stat"
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

func HttpGet(url string) (body []byte, statusCode int, status stat.Status) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("# Http Get - Panic occurred for '%s': %s\n", url, r)
		}
	}()

	for range only.Once {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("invalid URL '%s'", url),
				Help:    "if you provided this URL please correct it, otherwise " + stat.ContactSupportHelp(),
			})
			break
		}
		//request.SetBasicAuth(_api.User, _api.Password)
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		cli := &http.Client{}
		response, err := cli.Do(request)
		statusCode = response.StatusCode
		data := map[string]interface{}{
			"status_code": statusCode,
		}
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("http status code %d; unable to retrieve '%s': ",
					statusCode,
					url,
				),
				Help: "check your Internet connection",
			})
			break
		}
		defer CloseResponseBody(response)
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("http status code %d; unable to retrieve '%s': ",
					statusCode,
					url,
				),
				Help: "check your Internet connection",
			})
			break
		}
		data["response_body"] = body
		switch statusCode {
		case 200:
			if err != nil {
				status = stat.NewFailedStatus(&stat.Args{
					Error: err,
					Message: fmt.Sprintf("http status code 200 but no content returned for '%s'",
						url,
					),
					Help: stat.ContactSupportHelp(),
					Data: data,
				})
				break
			}

		case 401:
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("invalid credentials provided for '%s'",
					url,
				),
				Help: stat.ContactSupportHelp(),
				Data: data,
			})
			break

		case 403:
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("permission denied for '%s'",
					url,
				),
				Help: stat.ContactSupportHelp(),
				Data: data,
			})
			break

		default:
			status = stat.NewFailedStatus(&stat.Args{
				Error: err,
				Message: fmt.Sprintf("permission denied for '%s'",
					url,
				),
				Help: stat.ContactSupportHelp(),
				Data: data,
			})
			break
		}

	}
	return body, statusCode, status
}
