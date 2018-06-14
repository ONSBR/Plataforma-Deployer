package container

import (
	"context"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Container struct {
	Name  string
	ID    string
	Image string
}

//GetImages from docker
func GetImages() ([]Container, *exceptions.Exception) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	_containers := make([]Container, 0, 0)
	for _, container := range containers {
		_containers = append(_containers, Container{ID: container.ID, Name: container.Names[0], Image: container.Image})
	}
	return _containers, nil
}

func BuildImage() *exceptions.Exception {
	cli, err := client.NewEnvClient()
	if err != nil {
		return exceptions.NewComponentException(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	_containers := make([]Container, 0, 0)
	for _, container := range containers {
		_containers = append(_containers, Container{ID: container.ID, Name: container.Names[0], Image: container.Image})
	}
	return nil
}
