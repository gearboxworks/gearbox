package jsonapi

type DataObject struct {
	Data *ResourceObject `json:"data"`
}

func NewData(ro *ResourceObject) *DataObject {
	return &DataObject{
		Data: ro,
	}
}
