package dockerhub

type Repositories []*Respository

type Respository struct {
	User           string `json:"user"`
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	RepositoryType string `json:"repository_type"`
	Status         int    `json:"status"`
	Description    string `json:"description"`
	IsPrivate      bool   `json:"is_private"`
	IsAutomated    bool   `json:"is_automated"`
	CanEdit        bool   `json:"can_edit"`
	StarCount      int    `json:"star_count"`
	PullCount      int    `json:"pull_count"`
	LastUpdated    string `json:"last_updated"`
	IsMigrated     bool   `json:"is_migrated"`
}
