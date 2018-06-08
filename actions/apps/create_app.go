package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateApp install some app on APICore
func CreateApp(app *models.App) *exceptions.Exception {
	list := make([]models.App, 1)
	ex := apicore.FindByID(app.Metadata.Type, app.ID, &list)
	if ex != nil {
		return ex
	}
	if len(list) > 0 {
		return exceptions.NewInvalidArgumentException(fmt.Errorf("The app %s already exist", app.Name))
	}
	return apicore.PersistOne(app)
}
