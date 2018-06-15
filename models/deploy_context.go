package models

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	log "github.com/sirupsen/logrus"
)

//DeployContext is the entity to manage all deploy steps
type DeployContext struct {
	Info     *Deploy
	RootPath string
	Version  string
}

//GetImageTag returns docker image name pattern
func (context *DeployContext) GetImageTag() string {
	return fmt.Sprintf("%s/%s:%s", context.Info.App.SystemName, context.Info.Name, context.Version)
}

//GetDeployPath returns the location where sourcecode will be cloned
func (context *DeployContext) GetDeployPath() string {
	return fmt.Sprintf("%s/%s", env.GetDeploysPath(), context.Info.App.SystemName)
}

//Clone sourcecode from git
func (context *DeployContext) Clone() *exceptions.Exception {
	deployPath := context.GetDeployPath()
	url := fmt.Sprintf("%s/%s/%s", env.GetGitServerReposPath(), context.Info.App.SystemName, context.Info.Name)
	log.Info(fmt.Sprintf("Clonning code from %s to %s", url, deployPath))
	if ex := git.CloneRepo(deployPath, url, "master"); ex != nil {
		return ex
	}
	return nil
}

//Deploy register the function that will build the app and wrap with clone and cleanup procedures
func (context *DeployContext) Deploy(builder func(*DeployContext) *exceptions.Exception) *exceptions.Exception {
	if ex := context.Clone(); ex != nil {
		return ex
	}
	if ex := builder(context); ex != nil {
		return ex
	}
	if ex := context.Cleanup(); ex != nil {
		return ex
	}
	return nil
}

//Cleanup clear artifact deploy folder
func (context *DeployContext) Cleanup() *exceptions.Exception {
	deployPath := context.GetDeployPath()
	log.Info("Cleaning artifact folder")
	if err := os.RemoveAll(deployPath); err != nil {
		return exceptions.NewComponentException(err)
	}
	return nil
}
