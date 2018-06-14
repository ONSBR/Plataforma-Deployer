package main

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/actions/apps"
	"github.com/ONSBR/Plataforma-Deployer/api"
	"github.com/ONSBR/Plataforma-Deployer/container"
	log "github.com/sirupsen/logrus"
)

func init() {
	images := container.Images()
	fmt.Println(images)
	for _, image := range images {
		fmt.Println(image)
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	apps.RunDeployWorkers(1)
	api.Run()
}
