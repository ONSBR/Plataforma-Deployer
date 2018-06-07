package models

//BaseModel should be extended by all platform entities
type BaseModel struct {
	ID       string   `json:"id"`
	Metadata Metadata `json:"_metadata"`
}

//Metadata is the metadata for platform entities
type Metadata struct {
	Type        string `json:"type,omitempty"`
	ChangeTrack string `json:"changeTrack,omitempty"`
}
