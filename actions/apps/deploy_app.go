package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

var chProcessApps chan *models.App
var chDomainApps chan *models.App
var chPresentationApps chan *models.App

//RunDeployWorkers starts all listeners to deploy execution
func RunDeployWorkers(max int) {
	chProcessApps = make(chan *models.App)
	chDomainApps = make(chan *models.App)
	chPresentationApps = make(chan *models.App)

	for i := 0; i < max; i++ {
		go deployProcessAppWorker(chProcessApps)
		go deployDomainAppWorker(chDomainApps)
		go deployPresentationAppWorker(chPresentationApps)
	}

}

//DeployApp at platform based on app type
func DeployApp(app *models.App) *exceptions.Exception {
	switch app.Type {
	case "process":
		chProcessApps <- app
	case "domain":
		chDomainApps <- app
	case "presentation":
		chPresentationApps <- app
	}
	return nil
}
