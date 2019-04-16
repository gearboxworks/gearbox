package dockerhub

type ContainerNames []ContainerName

type ContainerName string

type ContainerMap map[ContainerName]*Container

type Containers []*Container

type Container struct {
	Id          int           `json:"id"`
	Name        ContainerName `json:"name"`
	FullSize    int           `json:"full_size"`
	Repository  int           `json:"repository"`
	Creator     int           `json:"creator"`
	LastUpdater int           `json:"last_updater"`
	LastUpdated string        `json:"last_updated"`
	ImageId     string        `json:"image_id"`
	V2          bool          `json:"v2"`
	//Images []*ContainerImage `json:"images"`
}

func RemoveContainer(i int, cs Containers) Containers {
	cs[i] = cs[len(cs)-1]
	return cs[:len(cs)-1]
}
