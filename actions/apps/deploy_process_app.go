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
	/*
		1) Faz o clone do repositório numa pasta temporaria
		2) Executa o docker build a partir do Dockerfile
		3) Grava as informações de versão do deploy na APICore e altera o status do deploy para deploying
		4) Grava as informações de operation na APICore (vide cli atual)
	*/
	return nil
}
