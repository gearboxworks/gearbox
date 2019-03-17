package api

type ListItemResponseMap map[string]*ListItemResponse

type ListItemResponse struct {
	Links Links       `json:"links"`
	Data  interface{} `json:"data"`
}

func NewListItemResponse(link UriTemplate, data interface{}) *ListItemResponse {
	return &ListItemResponse{
		Links: Links{
			SelfResource: link,
		},
		Data: data,
	}
}
