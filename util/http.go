package util

import (
	"fmt"
	"gearbox/help"
	"github.com/gearboxworks/go-status/only"

	"github.com/gearboxworks/go-status"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CloseResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Printf(
			"Could not close response body from HttpRequest: %s\n",
			err.Error(),
		)
	}
}

type HttpArgs struct {
	Method      string
	Body        io.Reader
	ContentType string
	Timeout     *time.Duration
}

func HttpRequest(url string, args ...*HttpArgs) (body []byte, statusCode int, sts status.Status) {

	var _args HttpArgs
	if len(args) == 0 {
		_args = HttpArgs{}
	} else {
		_args = *args[0]
	}
	if _args.Method == "" {
		_args.Method = "GET"
	}
	if _args.ContentType == "" {
		_args.ContentType = "application/json; charset=utf-8"
	}
	if _args.Timeout == nil {
		var d time.Duration
		d = time.Second * 3
		_args.Timeout = &d
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("# Http Get - Panic occurred for '%s': %s\n", url, r)
		}
	}()

	for range only.Once {
		request, err := http.NewRequest(_args.Method, url, _args.Body)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("invalid URL '%s'", url),
				Help:    fmt.Sprintf("if you provided this URL please correct it, otherwise %s", help.ContactSupportHelp()),
			})
			break
		}

		//request.SetBasicAuth(_api.User, _api.Password)
		request.Header.Set("Content-Type", _args.ContentType)
		cli := &http.Client{
			Timeout: *_args.Timeout,
		}
		response, err := cli.Do(request)
		if response != nil {
			statusCode = response.StatusCode
		} else {
			statusCode = http.StatusInternalServerError
		}
		data := map[string]interface{}{
			"status_code": statusCode,
		}
		if err != nil {
			sts = status.Wrap(err, &status.Args{
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
			sts = status.Wrap(err, &status.Args{
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
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("http status code 200 but no content returned for '%s'",
					url,
				),
				Help: help.ContactSupportHelp(),
				Data: data,
			})
			break

		case 401:
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("invalid credentials provided for '%s'",
					url,
				),
				Help: help.ContactSupportHelp(),
				Data: data,
			})
			break

		case 403:
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("permission denied for '%s'",
					url,
				),
				Help: help.ContactSupportHelp(),
				Data: data,
			})
			break

		default:
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("permission denied for '%s'",
					url,
				),
				Help: help.ContactSupportHelp(),
				Data: data,
			})
			break
		}

	}
	return body, statusCode, sts
}
