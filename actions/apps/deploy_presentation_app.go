package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployPresentationAppWorker(queue chan *models.DeployContext) {
	for deploy := range queue {
		log.Info(fmt.Sprintf("Deploying presentation app %s", deploy.Info.Name))
	}
}
