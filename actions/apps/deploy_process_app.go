package apps

import (
	"fmt"
	"io/ioutil"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployProcessAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		if ex := context.Deploy(doProcessDeploy); ex != nil {
			log.Error(ex.Message)
		} else {
			log.Info("Finished Deploy")
		}
	}
}

func doProcessDeploy(context *models.DeployContext) *exceptions.Exception {
	if files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", context.GetDeployPath(), context.Info.Name)); err != nil {
		return exceptions.NewComponentException(err)
	} else {
		for _, f := range files {
			log.Info(f.Name())
		}
	}

	return nil
}
