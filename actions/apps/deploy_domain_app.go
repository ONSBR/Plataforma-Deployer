package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployDomainAppWorker(queue chan *models.Deploy) {
	for deploy := range queue {
		log.Info(fmt.Sprintf("Deploying domain app %s", deploy.Name))
	}
}
