package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployPresentationAppWorker(queue chan *models.Deploy) {
	for deploy := range queue {
		log.Info(fmt.Sprintf("Deploying presentation app %s", deploy.Name))
	}
}
