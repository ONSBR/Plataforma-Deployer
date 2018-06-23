package models

import uuid "github.com/satori/go.uuid"

//Deploy is the entity to manage deploys on apicore
type Deploy struct {
	BaseModel
	SystemID    string `json:"systemId,omitempty"`
	ProcessID   string `json:"processId,omitempty"`
	Version     string `json:"version,omitempty"`
	Status      string `json:"status,omitempty"`
	Name        string `json:"name,omitempty"`
	Image       string `json:"image,omitempty"`
	ContainerID string `json:"containerId,omitempty"`
	App         *App   `json:"-"`
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
