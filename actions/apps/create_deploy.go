package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateDeploy at apicore
func CreateDeploy(processID string) *exceptions.Exception {
	if app, ex := FindAppByID(processID); ex != nil {
		return ex
	} else {
		deploy := models.NewDeploy()
		deploy.Name = app.Name
		deploy.ProcessID = app.ID
		deploy.SystemID = app.SystemID
		if ex := apicore.PersistOne(deploy); ex != nil {
			return ex
		}
		//Async deploy process
		go DeployApp(app)
	}
	return nil
}
