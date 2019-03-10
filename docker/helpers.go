package docker

import "fmt"

const (
	dockerHubEndpoint = "https://hub.docker.com"
	dockerHubVerion   = "v2"
)

var (
	URL = fmt.Sprintf("%s/%s", dockerHubEndpoint, dockerHubVerion)
)