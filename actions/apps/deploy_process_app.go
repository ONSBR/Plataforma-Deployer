package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"

	"github.com/ONSBR/Plataforma-Deployer/models"
	log "github.com/sirupsen/logrus"
)

func deployProcessAppWorker(queue chan *models.Deploy) {
	for deploy := range queue {
		log.Info(fmt.Sprintf("Deploying process app %s", deploy.App.Name))
		context := new(models.DeployContext)
		context.Info = deploy
		ex := doProcessDeploy(context)
		if ex != nil {
			//muda o status na apicore para aborted
			log.Error(ex.Message)
		}
	}
}

func doProcessDeploy(context *models.DeployContext) *exceptions.Exception {
	deployPath := fmt.Sprintf("%s/%s", env.GetDeploysPath(), context.Info.SystemID)
	url := env.GetSSHRemoteURL(context.Info.App.SystemName, context.Info.App.Name)
	git.CloneRepo(deployPath, url, "master")

	/*
		1) Faz o clone do repositório numa pasta temporaria
		2) Executa o docker build a partir do Dockerfile
		3) Grava as informações de versão do deploy na APICore e altera o status do deploy para deploying
		4) Grava as informações de operation na APICore (vide cli atual)
	*/
	return nil
}
