package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
)

func deployPresentationAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		context.Start(doPresentationDeploy)
	}
}

func doPresentationDeploy(context *models.DeployContext) error {
	context.RemoveContainer(context.GetContainerName())
	return nil
}
