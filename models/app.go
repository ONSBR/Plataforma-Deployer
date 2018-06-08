package models

//App is the model to represent all kind of platform apps
type App struct {
	BaseModel
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Type        string `json:"type"`
	SystemID    string `json:"systemId"`
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
