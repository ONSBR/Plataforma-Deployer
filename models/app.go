package models

//App is the model to represent all kind of platform apps
type App struct {
	BaseModel
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	SystemID    string `json:"systemId,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        uint   `json:"port,omitempty"`
	SystemName  string `json:"-"`
}

//IsDomain check if app is a domain app
func (app App) IsDomain() bool {
	return app.Type == "domain"
}

//IsProcess check if app is a process app
func (app App) IsProcess() bool {
	return app.Type == "process"
}

//IsPresentation check if app is a presentation app
func (app App) IsPresentation() bool {
	return app.Type == "presentation"
}

//NewApp builds a new App
func NewApp() *App {
	app := new(App)
	app.BaseModel.Metadata = Metadata{
		ChangeTrack: "create",
		Type:        "installedApp",
	}
	return app
}
