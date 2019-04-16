package ja

type JsonApi struct {
	Version Version `json:"version,omitempty"`
	MetaMap `json:"meta,omitempty"`
}
