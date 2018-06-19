package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/eventmanager"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployProcessAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		ex := context.Deploy(doProcessDeploy)
		if ex != nil {
			if ex := context.UpdateDeployStatus("error"); ex != nil {
				log.Error(ex)
			}
			log.Error(ex.Message)
		} else {
			log.Info("Finished Deploy")
			if ex := context.UpdateDeployStatus("success"); ex != nil {
				log.Error(ex)
			}
		}

		evt := eventmanager.Event{
			Name:    "system.deploy.finished",
			Payload: context.GetSummary(),
		}
		if ex := eventmanager.Push(&evt); ex != nil {
			log.Error(ex)
		}

	}
}

func doProcessDeploy(context *models.DeployContext) *exceptions.Exception {
	return nil
}
