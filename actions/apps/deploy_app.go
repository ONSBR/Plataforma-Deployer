package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

var chProcessApps chan *models.DeployContext
var chDomainApps chan *models.DeployContext
var chPresentationApps chan *models.DeployContext

//RunDeployWorkers starts all listeners to deploy execution
func RunDeployWorkers(max int) {
	chProcessApps = make(chan *models.DeployContext)
	chDomainApps = make(chan *models.DeployContext)
	chPresentationApps = make(chan *models.DeployContext)

	for i := 0; i < max; i++ {
		go deployProcessAppWorker(chProcessApps)
		go deployDomainAppWorker(chDomainApps)
		go deployPresentationAppWorker(chPresentationApps)
	}

}

//DeployApp at platform based on app type
func DeployApp(deploy *models.Deploy) *exceptions.Exception {
	context := new(models.DeployContext)
	context.Info = deploy
	switch deploy.App.Type {
	case "process":
		chProcessApps <- context
	case "domain":
		chDomainApps <- context
	case "presentation":
		chPresentationApps <- context
	}
	return nil
}
