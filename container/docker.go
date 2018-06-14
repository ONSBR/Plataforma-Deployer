package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Container struct {
	Name  string
	ID    string
	Image string
}

func Images() []Container {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	_containers := make([]Container, 0, 0)
	for _, container := range containers {
		_containers = append(_containers, Container{ID: container.ID, Name: container.Names[0], Image: container.Image})
	}
	return _containers
}
