package docker

import (
	"encoding/json"
	"fmt"
	"gearbox/util"
	"net/http"
)

func GetBuildHistory(a *Auth, user, repo string) (bh BuildHistory) {
	uri := fmt.Sprintf("/repositories/%s/%s/buildhistory/", user, repo)
	req, err := http.NewRequest("GET", URL+uri, nil)
	resource.Header.Set("Content-Type", "application/json")
	resource.Header.Set("Authorization", fmt.Sprintf("JWT %s", a.GetToken()))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer util.CloseResponseBody(resp)
	err = json.NewDecoder(resp.Body).Decode(&bh)
	if err != nil {
		return
	}
	return
}
