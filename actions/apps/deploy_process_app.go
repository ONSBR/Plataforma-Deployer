package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployProcessAppWorker(queue chan *models.App) {
	for app := range queue {
		log.Info(fmt.Sprintf("Deploying process app %s", app.Name))
	}
}
