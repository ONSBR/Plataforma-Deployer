package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
)

func deployProcessAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		context.Start(doProcessDeploy)
	}
}

func doProcessDeploy(context *models.DeployContext) error {
	return nil
}
