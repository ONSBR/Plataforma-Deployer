package models

type DependencyDomain struct {
	BaseModel
	ProcessID string `json:"processId"`
	SystemID  string `json:"systemId"`
	AppName   string `json:"name"`
	Entity    string `json:"entity"`
	Filter    string `json:"filter"`
	Version   string `json:"version"`
}

func NewDependencyDomain() *DependencyDomain {
	dd := new(DependencyDomain)
	dd.Metadata = Metadata{
		Type: "dependencyDomain",
	}
	return dd
}
