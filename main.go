package main

import (
	"os"

	"github.com/ONSBR/Plataforma-Deployer/actions/apps"
	"github.com/ONSBR/Plataforma-Deployer/api"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	apps.RunDeployWorkers(1)
	api.Run()
}
