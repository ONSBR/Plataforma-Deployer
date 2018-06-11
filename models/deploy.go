package models

//Deploy is the entity to manage deploys on apicore
type Deploy struct {
	Metadata
	SystemID  string `json:"systemId"`
	ProcessID string `json:"processId"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Name      string `json:"name"`
}

//NewDeploy creates a new deploy pointer
func NewDeploy() *Deploy {
	d := new(Deploy)
	d.Metadata = Metadata{
		ChangeTrack: "create",
		Type:        "deploy",
	}
	d.Status = "pending"
	return d
}
