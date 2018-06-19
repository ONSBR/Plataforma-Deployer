package models

type AppMap map[string]*EntityMap

type EntityMap struct {
	Model   string
	Fields  map[string]MapField
	Filters map[string]string
}

type MapField struct {
	Column string
}

func NewAppMap() AppMap {
	appMap := new(AppMap)
	return *appMap
}

type ApiCoreMap struct {
	BaseModel
	Name      string `json:"name"`
	SystemID  string `json:"systemId"`
	ProcessID string `json:"processId"`
	Content   string `json:"content"`
}

func NewApiCoreMap() *ApiCoreMap {
	m := new(ApiCoreMap)
	m.Metadata = Metadata{
		Type: "map",
	}
	return m
}
