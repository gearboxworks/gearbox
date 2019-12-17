package docker

type BuildHistoryItem struct {
	Id        int    `json:"id"`
	Status    int    `json:"status"`
	TagName   string `json:"dockertag_name"`
	BuildCode string `json:"build_code"`
}

type BuildHistory struct {
	Count int                `json:"count"`
	Itens []BuildHistoryItem `json:"results"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	Token string `json:"token"`
}