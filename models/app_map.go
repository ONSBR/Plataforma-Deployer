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
