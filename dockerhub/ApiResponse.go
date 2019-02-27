package dockerhub

type ApiResponse struct {
	Count    int
	Next     string
	Previous string
	Results  interface{}
}
