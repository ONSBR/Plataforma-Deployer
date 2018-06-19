package models

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/PMoneda/whaler"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/google/uuid"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	yaml "gopkg.in/yaml.v2"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	log "github.com/sirupsen/logrus"
)

//DeployContext is the entity to manage all deploy steps
type DeployContext struct {
	Info        *Deploy
	RootPath    string
	Version     string
	Metadata    *AppMetadata
	Map         AppMap
	MapName     string
	MapContent  string
	ContainerID string
	Error       *exceptions.Exception
}

type DeploySummary struct {
	ID        string `json:"deployId"`
	ProcessID string `json:"processId"`
	SystemID  string `json:"systemId"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

//GetDockerfilePath returns a path to app Dockerfile
func (context *DeployContext) GetDockerfilePath() string {
	return fmt.Sprintf("%s/Dockerfile", context.RootPath)
}

func (context *DeployContext) GetSummary() DeploySummary {
	summary := DeploySummary{
		ID:        context.Info.ID,
		ProcessID: context.Info.ProcessID,
		SystemID:  context.Info.SystemID,
	}
	status := "success"
	if context.Error != nil {
		status = "error"
		summary.Error = context.Error.Message
	}
	summary.Status = status
	return summary
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
	_version, err := uuid.NewUUID()
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	context.Version = _version.String()
	context.RootPath = fmt.Sprintf("%s/%s", context.GetDeployPath(), context.Info.Name)

	if ex := context.Clone(); ex != nil {
		context.Error = ex
		return ex
	}
	if !context.Info.App.IsDomain() {
		if ex := context.SaveAppMap(); ex != nil {
			context.Error = ex
			return ex
		}
		if ex := context.SaveDependencyDomain(); ex != nil {
			context.Error = ex
			return ex
		}
		if ex := context.SaveMetadata(); ex != nil {
			context.Error = ex
			return ex
		}
	}

	if ex := context.BuildImage(); ex != nil {
		context.Error = ex
		return ex
	}

	if !context.Info.App.IsProcess() {
		if ex := context.StartApp(); ex != nil {
			context.Error = ex
			return ex
		}
	}

	if ex := builder(context); ex != nil {
		context.Error = ex
		return ex
	}
	if ex := context.Cleanup(); ex != nil {
		context.Error = ex
		return ex
	}
	return nil
}

func (context *DeployContext) UpdateDeployStatus(status string) *exceptions.Exception {
	dep := NewDeploy()
	dep.ID = context.Info.ID
	dep.Status = status
	dep.Metadata.ChangeTrack = "update"
	if ex := apicore.PersistOne(dep); ex != nil {
		context.Error = ex
		return ex
	}
	return nil
}

func (context *DeployContext) BuildImage() *exceptions.Exception {
	build := func(worker string) *exceptions.Exception {
		cnf := whaler.BuildImageConfig{
			PathContext: context.RootPath,
			Tag:         context.GetImageName(""),
		}
		_, err := whaler.BuildImageWithDockerfile(cnf)
		if err != nil {
			return exceptions.NewComponentException(err)
		}
		return nil
	}
	if ex := build(""); ex != nil {
		return ex
	}
	if context.Info.App.IsDomain() {
		return build("worker")
	}
	return nil
}

func (context *DeployContext) StartApp() *exceptions.Exception {
	externalPort := "8087"
	if context.Info.App.IsPresentation() {
		externalPort = "8088"
	}
	rand.Seed(int64(time.Now().Nanosecond()))

	debugPort := fmt.Sprintf("%d%d%d%d", 7, rand.Intn(9), rand.Intn(9), rand.Intn(9))

	cnf := whaler.CreateContainerConfig{
		Image:       context.GetImageName(""),
		NetworkName: "plataforma_network",
		Name:        context.GetContainerName(),
		Env:         context.GetEnvVars(),
		Ports:       []string{fmt.Sprintf("%s:%s", externalPort, context.GetPort()), fmt.Sprintf("%s:%s", debugPort, context.GetDebugPort())},
	}
	log.Info(cnf)
	id, err := whaler.CreateContainer(cnf)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	context.ContainerID = id
	if err := whaler.StartContainer(id); err != nil {
		return exceptions.NewComponentException(err)
	}
	return nil
}

func (context *DeployContext) SaveDependencyDomain() *exceptions.Exception {
	log.Info("saving dependency domain")
	list := make([]*DependencyDomain, 0)
	if ex := apicore.FindByProcessID(NewDependencyDomain().Metadata.Type, context.Info.ProcessID, &list); ex != nil {
		return ex
	}
	for k, v := range context.Map {
		for _, f := range v.Filters {
			dd := NewDependencyDomain()
			dd.Metadata.ChangeTrack = "create"
			dd.AppName = context.Info.App.Name
			dd.Entity = k
			dd.ProcessID = context.Info.ProcessID
			dd.SystemID = context.Info.SystemID
			dd.Version = context.Version
			dd.Filter = f
			list = append(list, dd)
		}
	}
	return apicore.Persist(list)
}

func (context *DeployContext) GetEnvVars() []string {
	if context.Info.App.IsPresentation() {
		return []string{"API_MODE=true", fmt.Sprintf("SYSTEM_ID=%s", context.Info.SystemID)}
	}
	return nil
}

func (context *DeployContext) GetPort() string {
	if context.Info.App.IsDomain() {
		return "9110"
	}
	return "8088"
}

func (context *DeployContext) GetDebugPort() string {
	return "9229" //node debug default
}

//GetMetadata returns a metadata configuration app
func (context *DeployContext) GetMetadata() (*AppMetadata, *exceptions.Exception) {
	meta := NewAppMetadata()
	path := fmt.Sprintf("%s/metadados", context.RootPath)
	data, _, ex := context.readFirstFileInDir(path)
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

func (context *DeployContext) SaveMetadata() *exceptions.Exception {
	meta, ex := context.GetMetadata()
	if ex != nil {
		return ex
	}
	operations := make([]*OperationCore, len(meta.Operations))
	i := 0
	list := make([]*OperationCore, 0)
	if ex := apicore.FindByProcessID("operation", context.Info.ProcessID, &list); ex != nil {
		return ex
	}
	setID := func(op *OperationCore) *exceptions.Exception {
		for _, o := range list {
			if o.EventIn == op.EventIn {
				if o.ProcessID != op.ProcessID {
					return exceptions.NewInvalidArgumentException(fmt.Errorf("Conflict event in operation %s is already mapped to app %s", op.EventIn, o.Name))
				}
				op.ID = o.ID
				op.Metadata.ChangeTrack = "update"
				return nil
			}
		}
		return nil
	}
	for _, op := range meta.Operations {
		coreOp := NewOperationCore()
		coreOp.Metadata.ChangeTrack = "create"
		coreOp.EventIn = op.Event
		coreOp.EventOut = context.Info.Name + ".done"
		coreOp.Name = op.Name
		coreOp.Commit = op.Commit
		coreOp.ProcessID = context.Info.ProcessID
		coreOp.SystemID = context.Info.SystemID
		coreOp.Version = context.Version
		coreOp.Image = context.GetImageName("")
		if ex := setID(coreOp); ex != nil {
			return ex
		}
		operations[i] = coreOp
		i++
	}
	return apicore.Persist(operations)
}

func (context *DeployContext) GetImageName(worker string) string {
	if worker != "" {
		return fmt.Sprintf("registry:5000/%s_%s:%s", context.Info.App.Name, worker, context.Version)
	}
	return fmt.Sprintf("registry:5000/%s:%s", context.Info.App.Name, context.Version)
}

func (context *DeployContext) GetContainerName() string {
	name := context.Info.App.Name
	if context.Info.App.SystemName != "plataforma" {
		name = fmt.Sprintf("%s-%s", context.Info.App.SystemName, name)
	}
	return name
}

//SaveAppMap saves application map to apicore
func (context *DeployContext) SaveAppMap() *exceptions.Exception {
	_, ex := context.GetAppMap()
	if ex != nil {
		return ex
	}
	apiCoreMap := NewApiCoreMap()
	apiCoreMap.ProcessID = context.Info.ProcessID
	apiCoreMap.SystemID = context.Info.SystemID
	apiCoreMap.Name = context.MapName
	apiCoreMap.Content = context.MapContent
	existingMap := make([]*ApiCoreMap, 0)
	if ex := apicore.FindByProcessID("map", context.Info.ProcessID, &existingMap); ex != nil {
		return ex
	}
	if len(existingMap) == 0 {
		apiCoreMap.Metadata.ChangeTrack = "create"
	} else {
		apiCoreMap.ID = existingMap[0].ID
		apiCoreMap.Metadata.ChangeTrack = "update"
	}
	return apicore.PersistOne(apiCoreMap)
}

//GetAppMap returns a domain map defined by app
func (context *DeployContext) GetAppMap() (AppMap, *exceptions.Exception) {
	mapApp := NewAppMap()
	path := fmt.Sprintf("%s/mapa", context.RootPath)
	log.Info(path)
	data, fileName, ex := context.readFirstFileInDir(path)
	if ex != nil {
		return nil, ex
	}
	err := yaml.Unmarshal(data, &mapApp)
	if err != nil {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("Invalid yaml format: %s", err.Error()))
	}
	context.Map = mapApp
	context.MapContent = string(data)
	parts := strings.Split(fileName, ".")
	context.MapName = parts[0]
	return context.Map, nil
}

func (context *DeployContext) readFirstFileInDir(path string) ([]byte, string, *exceptions.Exception) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, "", exceptions.NewComponentException(err)
	}
	if len(files) == 0 {
		return nil, "", exceptions.NewInvalidArgumentException(fmt.Errorf("no file found in %s", path))
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, files[0].Name()))
	if err != nil {
		return nil, "", exceptions.NewComponentException(err)
	}
	return data, files[0].Name(), nil
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
