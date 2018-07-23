package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateDeploy at apicore
func CreateDeploy(processID string) error {
	if app, ex := FindAppByID(processID); ex != nil {
		return ex
	} else {
		deploy := models.NewDeploy()
		deploy.Name = app.Name
		deploy.ProcessID = app.ID
		deploy.SystemID = app.SystemID
		if sol, ex := FindSolutionByID(deploy.SystemID); ex != nil {
			return ex
		} else {
			app.SystemName = sol.Name
		}
		if ex := apicore.PersistOne(deploy); ex != nil {
			return ex
		}
		//Async deploy process
		deploy.App = app
		go DeployApp(deploy)
	}
	return nil
}
