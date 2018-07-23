package models

type AppMap map[string]*EntityMap

type EntityMap struct {
	Model   string
	Fields  map[string]MapField
	Filters map[string]string
}

func (m AppMap) HasFilters() bool {
	for _, v := range m {
		if len(v.Filters) > 0 {
			return true
		}
	}
	return false
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
