package solution

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateSolution on platform
func CreateSolution(solution *models.Solution) *exceptions.Exception {
	ex := createSolutionOnAPICore(solution)
	if ex != nil {
		return ex
	}
	return createRootFolderForSolutionOnGit(solution)
}

func createSolutionOnAPICore(solution *models.Solution) *exceptions.Exception {
	sol, ex := FindSolutionByID(solution.ID)
	if ex != nil {
		return ex
	}
	if sol != nil {
		return exceptions.NewInvalidArgumentException(fmt.Errorf("solution %s already exists", solution.ID))
	}
	return apicore.PersistOne(solution)
}

func createRootFolderForSolutionOnGit(solution *models.Solution) *exceptions.Exception {
	path := fmt.Sprintf("%s/%s", env.GetGitServerReposPath(), solution.Name)
	err := os.Mkdir(path, 0777)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	return nil
}
