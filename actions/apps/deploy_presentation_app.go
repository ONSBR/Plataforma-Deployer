package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

func deployPresentationAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		context.Start(doPresentationDeploy)
	}
}

func doPresentationDeploy(context *models.DeployContext) *exceptions.Exception {
	context.RemoveContainer(context.GetContainerName())
	return nil
}
