package docker

import (
	"bytes"
	"encoding/json"
	"gearbox/util"
	"net/http"
)

func (a *Auth) GetToken() (l Login) {
	uri := "/users/login/"
	bff := new(bytes.Buffer)
	err := json.NewEncoder(bff).Encode(a)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", URL+uri, bff)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer util.CloseResponseBody(resp)
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		panic(err)
	}
	return
}

