package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

var chProcessApps chan *models.Deploy
var chDomainApps chan *models.Deploy
var chPresentationApps chan *models.Deploy

//RunDeployWorkers starts all listeners to deploy execution
func RunDeployWorkers(max int) {
	chProcessApps = make(chan *models.Deploy)
	chDomainApps = make(chan *models.Deploy)
	chPresentationApps = make(chan *models.Deploy)

	for i := 0; i < max; i++ {
		go deployProcessAppWorker(chProcessApps)
		go deployDomainAppWorker(chDomainApps)
		go deployPresentationAppWorker(chPresentationApps)
	}

}

//DeployApp at platform based on app type
func DeployApp(deploy *models.Deploy) *exceptions.Exception {
	switch deploy.App.Type {
	case "process":
		chProcessApps <- deploy
	case "domain":
		chDomainApps <- deploy
	case "presentation":
		chPresentationApps <- deploy
	}
	return nil
}
