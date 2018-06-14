package main

import (
	"fmt"
	"os"

	"github.com/PMoneda/whaler"

	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	//apps.RunDeployWorkers(1)
	//api.Run()
	//"/home/philippe/installed_plataforma/Plataforma-Installer/Dockerfiles"
	str, err := whaler.BuildImageWithDockerfile(whaler.BuildImageConfig{
		Dockerfile:  "EventManager",
		PathContext: "/home/philippe/installed_plataforma/Plataforma-Installer/Dockerfiles",
		Tag:         "plataforma/meu-container:1.0",
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(str)
	id, err := whaler.CreateContainer(whaler.CreateContainerConfig{
		Image: "plataforma/meu-container:1.0",
		Name:  "meu-container",
		Ports: []string{"8081:8081"},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(id)

	err = whaler.StartContainer(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
