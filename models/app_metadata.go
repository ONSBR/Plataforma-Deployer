package models

//AppMetadata map apps metadata file configuration
type AppMetadata struct {
	Operations []*Operation
}

//Operation defines app binding configuration
type Operation struct {
	Name             string
	Event            string
	Commit           bool
	SkipReprocessing bool `yaml:"skip_reprocessing"`
}

//NewAppMetadata returns a new instance of AppMetadata
func NewAppMetadata() *AppMetadata {
	meta := new(AppMetadata)
	meta.Operations = make([]*Operation, 0)
	return meta
}
