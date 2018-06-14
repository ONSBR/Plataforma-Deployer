package container

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/PMoneda/whaler"
)

//BuildApp builds a docker image and container for a specific deploy context
func BuildApp(deploy models.DeployContext) (string, *exceptions.Exception) {
	_, err := whaler.BuildImageWithDockerfile(whaler.BuildImageConfig{
		PathContext: deploy.RootPath,
		Tag:         deploy.GetImageTag(),
	})
	if err != nil {
		return "", exceptions.NewComponentException(err)
	}
	id, err := whaler.CreateContainer(whaler.CreateContainerConfig{
		Image:       deploy.GetImageTag(),
		Name:        fmt.Sprintf("%s_%s", deploy.Solution.Name, deploy.Info.Name),
		NetworkName: "plataforma_network",
	})
	if err != nil {
		return "", exceptions.NewComponentException(err)
	}
	return id, nil
}
