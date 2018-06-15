package models

import uuid "github.com/satori/go.uuid"

//Deploy is the entity to manage deploys on apicore
type Deploy struct {
	BaseModel
	SystemID  string `json:"systemId"`
	ProcessID string `json:"processId"`
	Version   string `json:"version,omitempty"`
	Status    string `json:"status,omitempty"`
	Name      string `json:"name"`
	App       *App   `json:"-"`
}

//NewDeploy creates a new deploy pointer
func NewDeploy() *Deploy {
	d := new(Deploy)
	d.ID = uuid.NewV4().String()
	d.Metadata = Metadata{
		ChangeTrack: "create",
		Type:        "deploy",
	}
	d.Version = uuid.NewV4().String()
	d.Status = "pending"
	return d
}
