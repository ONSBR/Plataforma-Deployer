package main

import (
	"fmt"
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
	print()
	apps.RunDeployWorkers(1)
	api.Run()
}

func print() {
	fmt.Print(`

	$$$$$$$\                      $$\
	$$  __$$\                     $$ |
	$$ |  $$ | $$$$$$\   $$$$$$\  $$ | $$$$$$\  $$\   $$\  $$$$$$\   $$$$$$\
	$$ |  $$ |$$  __$$\ $$  __$$\ $$ |$$  __$$\ $$ |  $$ |$$  __$$\ $$  __$$\
	$$ |  $$ |$$$$$$$$ |$$ /  $$ |$$ |$$ /  $$ |$$ |  $$ |$$$$$$$$ |$$ |  \__|
	$$ |  $$ |$$   ____|$$ |  $$ |$$ |$$ |  $$ |$$ |  $$ |$$   ____|$$ |
	$$$$$$$  |\$$$$$$$\ $$$$$$$  |$$ |\$$$$$$  |\$$$$$$$ |\$$$$$$$\ $$ |
	\_______/  \_______|$$  ____/ \__| \______/  \____$$ | \_______|\__|
	                    $$ |                    $$\   $$ |
	                    $$ |                    \$$$$$$  |
	                    \__|                     \______/


	`)
}
