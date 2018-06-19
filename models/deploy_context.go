package models

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	yaml "gopkg.in/yaml.v2"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	log "github.com/sirupsen/logrus"
)

//DeployContext is the entity to manage all deploy steps
type DeployContext struct {
	Info     *Deploy
	RootPath string
	Version  string
	Metadata *AppMetadata
	Map      AppMap
}

//GetDockerfilePath returns a path to app Dockerfile
func (context *DeployContext) GetDockerfilePath() string {
	return fmt.Sprintf("%s/Dockerfile", context.RootPath)
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

//GetMetadata returns a metadata configuration app
func (context *DeployContext) GetMetadata() (*AppMetadata, *exceptions.Exception) {
	meta := NewAppMetadata()
	path := fmt.Sprintf("%s/metadados", context.RootPath)
	data, ex := context.readFirstFileInDir(path)
	if ex != nil {
		return nil, ex
	}
	err := yaml.Unmarshal(data, meta)
	if err != nil {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("Invalid yaml format: %s", err.Error()))
	}
	context.Metadata = meta
	return context.Metadata, nil
}

//GetAppMap returns a domain map defined by app
func (context *DeployContext) GetAppMap() (AppMap, *exceptions.Exception) {
	mapApp := NewAppMap()
	path := fmt.Sprintf("%s/mapa", context.RootPath)
	data, ex := context.readFirstFileInDir(path)
	if ex != nil {
		return nil, ex
	}
	err := yaml.Unmarshal(data, &mapApp)
	if err != nil {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("Invalid yaml format: %s", err.Error()))
	}
	context.Map = mapApp
	return context.Map, nil
}

func (context *DeployContext) readFirstFileInDir(path string) ([]byte, *exceptions.Exception) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	if len(files) == 0 {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("no file found in %s", path))
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, files[0].Name()))
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	return data, nil
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
