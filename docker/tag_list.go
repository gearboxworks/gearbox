package docker

type TagList []string

type TagListResponse struct {
	Name string `json:"name"`
	Tags TagList `json:"tags"`
}

