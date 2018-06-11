package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployPresentationAppWorker(queue chan *models.App) {
	for app := range queue {
		log.Info(fmt.Sprintf("Deploying presentation app %s", app.Name))
	}
}
