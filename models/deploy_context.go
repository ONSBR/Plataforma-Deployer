package models

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/PMoneda/whaler"

	"github.com/ONSBR/Plataforma-Deployer/builder"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Deployer/sdk/eventmanager"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/google/uuid"

	yaml "gopkg.in/yaml.v2"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	log "github.com/sirupsen/logrus"
)

//DeployContext is the entity to manage all deploy steps
type DeployContext struct {
	Info       *Deploy
	RootPath   string
	Version    string
	Metadata   *AppMetadata
	Map        AppMap
	MapName    string
	MapContent string
	Stack      string
	Error      error
}

type DeploySummary struct {
	ID        string `json:"deployId"`
	ProcessID string `json:"processId"`
	SystemID  string `json:"systemId"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

func (summary DeploySummary) ToEventPayload() map[string]interface{} {
	_map := make(map[string]interface{})
	_map["deployId"] = summary.ID
	_map["processId"] = summary.ProcessID
	_map["systemId"] = summary.SystemID
	_map["status"] = summary.Status
	if summary.Error != "" {
		_map["error"] = summary.Error
	}
	return _map
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
		summary.Error = context.Error.Error()
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

//GetWorkspace returns the location of current deploying app
func (context *DeployContext) GetWorkspace() string {
	return fmt.Sprintf("%s/%s", context.GetDeployPath(), context.Info.App.Name)
}

//Clone sourcecode from git
func (context *DeployContext) Clone() error {
	deployPath := context.GetDeployPath()
	url := fmt.Sprintf("%s/%s/%s", env.GetGitServerReposPath(), context.Info.App.SystemName, context.Info.Name)
	log.Info(fmt.Sprintf("Clonning code from %s to %s", url, deployPath))
	if ex := git.CloneRepo(deployPath, url, "master"); ex != nil {
		return ex
	}
	return nil
}

//RemoveContainer removes docker container
func (context *DeployContext) RemoveContainer(name string) error {
	if container, err := whaler.FindContainerByIdentifier("/"+name, true); err == nil {
		return whaler.RemoveContainer(container.ID, true)
	} else {
		return err
	}
}

func (context *DeployContext) Start(builder func(*DeployContext) error) {
	ex := context.Deploy(builder)
	if ex != nil {
		if ex := context.UpdateDeployStatus("error"); ex != nil {
			log.Error(ex)
		}
		log.Error(ex)
	} else {
		log.Info("Finished Deploy")
		if ex := context.UpdateDeployStatus("success"); ex != nil {
			log.Error(ex)
		}
	}
	evt := domain.Event{
		Name:    "system.deploy.finished",
		Payload: context.GetSummary().ToEventPayload(),
	}
	if ex := eventmanager.Push(&evt); ex != nil {
		log.Error(ex)
	}
}

//Deploy register the function that will build the app and wrap with clone and cleanup procedures
func (context *DeployContext) Deploy(builder func(*DeployContext) error) error {
	_version, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	context.Version = _version.String()
	context.RootPath = fmt.Sprintf("%s/%s", context.GetDeployPath(), context.Info.Name)
	if ex := context.Cleanup(); ex != nil {
		context.Error = ex
		return ex
	}
	if ex := context.Clone(); ex != nil {
		context.Error = ex
		return ex
	}
	if ex := builder(context); ex != nil {
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

	if ex := context.IdentifyStack(); ex != nil {
		context.Error = ex
		return ex
	}

	/*if ex := context.RunStackDefaultBuild(); ex != nil {
		context.Error = ex
		return ex
	}*/

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
	return nil
}

func (context *DeployContext) IdentifyStack() error {
	files, err := ioutil.ReadDir(context.GetWorkspace())
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.Name() == "package.json" {
			context.Stack = "node"
			return nil
		}
	}
	context.Stack = "dotnet"
	return nil
}

func (context *DeployContext) RunStackDefaultBuild() error {
	builder := builder.GetBuilder(context.Stack)
	if builder != nil {
		return builder.Build(context.GetWorkspace())
	}
	return nil
}

func (context *DeployContext) UpdateDeployStatus(status string) error {
	dep := NewDeploy()
	dep.ID = context.Info.ID
	dep.Status = status
	dep.Version = context.Version
	dep.Image = context.GetImageName("")
	dep.ContainerID = context.Info.ContainerID
	dep.Metadata.ChangeTrack = "update"
	if ex := apicore.PersistOne(dep); ex != nil {
		context.Error = ex
		return ex
	}
	return nil
}

func (context *DeployContext) BuildImage() error {
	build := func(worker string) error {
		cnf := whaler.BuildImageConfig{
			PathContext: context.RootPath,
			Tag:         context.GetImageName(worker),
		}
		_, err := whaler.BuildImageWithDockerfile(cnf)
		if err != nil {
			return err
		}
		if _, err := whaler.Publish(context.GetImageName(worker), "docker", "docker"); err != nil {
			return err
		}
		return nil
	}
	if ex := build(""); ex != nil {
		return ex
	}
	return nil
}

func (context *DeployContext) StartApp() error {
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
		return err
	}
	context.Info.ContainerID = id
	if err := whaler.StartContainer(id); err != nil {
		return err
	}
	return nil
}

func (context *DeployContext) SaveDependencyDomain() error {
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
	if context.Info.App.IsDomain() {
		return []string{"DOMAIN_API=1"}
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
func (context *DeployContext) GetMetadata() (*AppMetadata, error) {
	meta := NewAppMetadata()
	path := fmt.Sprintf("%s/metadados", context.RootPath)
	data, _, ex := context.readFirstFileInDir(path)
	if ex != nil {
		return nil, ex
	}
	err := yaml.Unmarshal(data, meta)
	if err != nil {
		return nil, fmt.Errorf("Invalid yaml format: %s", err.Error())
	}
	context.Metadata = meta
	return context.Metadata, nil
}

func (context *DeployContext) SaveMetadata() error {
	meta, ex := context.GetMetadata()
	if ex != nil {
		return ex
	}
	return context.PersistOperations(meta.Operations)
}

func (context *DeployContext) PersistOperations(ops []*Operation) error {
	operations := make([]*OperationCore, len(ops))
	for i, op := range ops {
		coreOp := NewOperationCore()
		coreOp.Metadata.ChangeTrack = "create"
		coreOp.EventIn = op.Event
		if strings.HasSuffix(op.Event, ".request") {
			coreOp.EventOut = strings.Replace(op.Event, "request", "done", 1)
		} else {
			coreOp.EventOut = op.Event + ".done"
		}
		coreOp.Name = op.Name
		coreOp.Commit = op.Commit
		coreOp.Reprocessable = context.Map.HasFilters() && !op.SkipReprocessing
		coreOp.ProcessID = context.Info.ProcessID
		coreOp.SystemID = context.Info.SystemID
		coreOp.Version = context.Version
		coreOp.Image = context.GetImageName("")
		operations[i] = coreOp
	}
	return apicore.Persist(operations)
}

func (context *DeployContext) GetImageName(worker string) string {
	if worker != "" {
		return fmt.Sprintf("localhost:5000/%s_%s:%s", context.Info.App.Name, worker, context.Version)
	}
	return fmt.Sprintf("localhost:5000/%s:%s", context.Info.App.Name, context.Version)
}

func (context *DeployContext) GetContainerName() string {
	name := context.Info.App.Name
	if context.Info.App.SystemName != "plataforma" {
		name = fmt.Sprintf("%s-%s", context.Info.App.SystemName, name)
	}
	return name
}

//SaveAppMap saves application map to apicore
func (context *DeployContext) SaveAppMap() error {
	_, ex := context.GetAppMap()
	if ex != nil {
		return ex
	}
	apiCoreMap := NewApiCoreMap()
	apiCoreMap.ProcessID = context.Info.ProcessID
	apiCoreMap.SystemID = context.Info.SystemID
	apiCoreMap.Name = context.MapName
	apiCoreMap.Version = context.Version
	apiCoreMap.Content = context.MapContent
	apiCoreMap.Metadata.ChangeTrack = "create"
	return apicore.PersistOne(apiCoreMap)
}

//GetAppMap returns a domain map defined by app
func (context *DeployContext) GetAppMap() (AppMap, error) {
	mapApp := NewAppMap()
	path := fmt.Sprintf("%s/mapa", context.RootPath)
	log.Info(path)
	data, fileName, ex := context.readFirstFileInDir(path)
	if ex != nil {
		return nil, ex
	}
	err := yaml.Unmarshal(data, &mapApp)
	if err != nil {
		return nil, fmt.Errorf("Invalid yaml format: %s", err.Error())
	}
	context.Map = mapApp
	context.MapContent = string(data)
	parts := strings.Split(fileName, ".")
	context.MapName = parts[0]
	return context.Map, nil
}

func (context *DeployContext) readFirstFileInDir(path string) ([]byte, string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, "", err
	}
	if len(files) == 0 {
		return nil, "", fmt.Errorf("no file found in %s", path)
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, files[0].Name()))
	if err != nil {
		return nil, "", err
	}
	return data, files[0].Name(), nil
}

//Cleanup clear artifact deploy folder
func (context *DeployContext) Cleanup() error {
	deployPath := context.GetDeployPath()
	log.Info("Cleaning artifact folder")
	if err := os.RemoveAll(deployPath); err != nil {
		return err
	}
	if context.Info.App.IsDomain() {
		if err := os.RemoveAll(fmt.Sprintf("%s/%s", env.GetDeploysPath(), context.Info.ID)); err != nil {
			return err
		}
	}
	return nil
}
