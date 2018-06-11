package apps

import (
	"fmt"
	"io/ioutil"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/dto"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateApp install some app on APICore
func CreateApp(app *models.App) (*dto.CreateAppResponse, *exceptions.Exception) {
	if app.ID == "" {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("id is required"))
	}

	if ex := checkIfAppExist(app); ex != nil {
		return nil, ex
	}
	solution, err := FindSolutionByID(app.SystemID)
	if err != nil {
		return nil, err
	}
	if err := installApp(app, solution); err != nil {
		return nil, err
	}
	resp := dto.CreateAppResponse{
		GitRemote: env.GetSSHRemoteURL(solution.Name, app.Name),
	}
	return &resp, nil
}

func createGitRepo(solution, name string) *exceptions.Exception {
	return git.CreateGitRepo(fmt.Sprintf("%s/%s/%s", env.GetGitServerReposPath(), solution, name))
}

func checkIfAppExist(app *models.App) *exceptions.Exception {
	list := make([]models.App, 1)
	ex := apicore.FindByID(app.Metadata.Type, app.ID, &list)
	if ex != nil {
		return ex
	}
	if len(list) > 0 {
		return exceptions.NewInvalidArgumentException(fmt.Errorf("The app %s already exist", app.Name))
	}
	return nil
}

func setGitHook(app *models.App, solution *models.Solution) *exceptions.Exception {
	fileName := fmt.Sprintf("%s/%s/%s/hooks/post-receive", env.GetGitServerReposPath(), solution.Name, app.Name)
	scheme := env.Get("DEPLOYER_SCHEME", "http")
	host := env.Get("DEPLOYER_HOST", "localhost")
	port := env.Get("DEPLOYER_PORT", "6970")
	ex := ioutil.WriteFile(fileName, []byte(fmt.Sprintf(`#!/bin/sh

exec curl -X POST %s://%s:%s/api/v1.0.0/app/%s/deploy -d ''
		`, scheme, host, port, app.ID)), 0777)

	if ex != nil {
		return exceptions.NewComponentException(ex)
	}
	return nil
}

func installApp(app *models.App, solution *models.Solution) *exceptions.Exception {
	if err := apicore.PersistOne(app); err != nil {
		return err
	}
	if err := createGitRepo(solution.Name, app.Name); err != nil {
		return err
	}
	if err := setGitHook(app, solution); err != nil {
		return err
	}
	return nil
}
