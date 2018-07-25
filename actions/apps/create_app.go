package apps

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/git"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/dto"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateApp install some app on APICore
func CreateApp(app *models.App) (*dto.CreateAppResponse, error) {
	if app.ID == "" {
		return nil, fmt.Errorf("id is required")
	}
	solution, err := FindSolutionByID(app.SystemID)
	if err != nil {
		return nil, err
	}
	if solution == nil {
		return nil, fmt.Errorf("solution not found")
	}
	if ex := checkIfAppExist(app); ex != nil {
		resp := dto.CreateAppResponse{
			GitRemote: env.GetSSHRemoteURL(solution.Name, app.Name),
		}
		return &resp, nil
	}
	app.SystemName = solution.Name
	if err := installApp(app, solution); err != nil {
		return nil, err
	}
	resp := dto.CreateAppResponse{
		GitRemote: env.GetSSHRemoteURL(solution.Name, app.Name),
	}
	return &resp, nil
}

func createGitRepo(solution, name string) error {
	return git.CreateGitRepo(fmt.Sprintf("%s/%s/%s", env.GetGitServerReposPath(), solution, name))
}

func checkIfAppExist(app *models.App) error {
	list := make([]models.App, 1)
	ex := apicore.FindByID(app.Metadata.Type, app.ID, &list)
	if ex != nil {
		return ex
	}
	if len(list) > 0 {
		return fmt.Errorf("The app %s already exist", app.Name)
	}
	return nil
}

func setGitHook(app *models.App, solution *models.Solution) error {
	fileName := fmt.Sprintf("%s/%s/%s/hooks/post-receive", env.GetGitServerReposPath(), solution.Name, app.Name)
	scheme := env.Get("DEPLOYER_SCHEME", "http")
	host := env.Get("DEPLOYER_HOST", "localhost")
	port := env.Get("DEPLOYER_PORT", "6970")
	ex := ioutil.WriteFile(fileName, []byte(fmt.Sprintf(`#!/bin/sh

exec curl -X POST %s://%s:%s/api/v1.0.0/app/%s/deploy -d ''
		`, scheme, host, port, app.ID)), 0777)

	if ex != nil {
		return ex
	}
	return nil
}

func installApp(app *models.App, solution *models.Solution) error {
	if app.IsDomain() {
		app.Host = app.SystemName + "-" + app.Name
		app.Port = 9110
	} else if app.IsPresentation() {
		app.Host = app.SystemName + "-" + app.Name
		app.Port = 8088
	}
	if err := apicore.PersistOne(app); err != nil {
		return err
	}
	if err := createGitRepo(solution.Name, app.Name); err != nil {
		app.Metadata.ChangeTrack = "destroy"
		if err1 := apicore.PersistOne(app); err1 != nil {
			err = fmt.Errorf("%s/%s", err, err1)
		}
		return err
	}
	if err := setGitHook(app, solution); err != nil {
		app.Metadata.ChangeTrack = "destroy"
		if err1 := apicore.PersistOne(app); err1 != nil {
			err = fmt.Errorf("%s/%s", err, err1)
		}
		if err2 := os.RemoveAll(fmt.Sprintf("%s/%s/%s", env.GetGitServerReposPath(), solution.Name, app.Name)); err2 != nil {
			err = fmt.Errorf("%s/%s", err, err2)
		}
		return err
	}
	return nil
}
