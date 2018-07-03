package apps

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateSolution on platform
func CreateSolution(solution *models.Solution) error {
	ex := createSolutionOnAPICore(solution)
	if ex != nil {
		return ex
	}
	if ex := createRootFolderForSolutionOnGit(solution); ex != nil {
		if ex1 := deleteSolutionOnAPICore(solution); ex1 != nil {
			ex = fmt.Errorf("%s/%s", ex, ex1)
		}
		return ex
	}
	return nil
}

func deleteSolutionOnAPICore(solution *models.Solution) error {
	solution.Metadata.ChangeTrack = "destroy"
	return apicore.PersistOne(solution)
}

func createSolutionOnAPICore(solution *models.Solution) error {
	sol, ex := FindSolutionByID(solution.ID)
	if ex != nil {
		return ex
	}
	if sol.Name != "" {
		return fmt.Errorf("solution %s already exists", solution.ID)
	}
	return apicore.PersistOne(solution)
}

func createRootFolderForSolutionOnGit(solution *models.Solution) error {
	path := fmt.Sprintf("%s/%s", env.GetGitServerReposPath(), solution.Name)
	err := os.Mkdir(path, 0777)
	if err != nil {
		return err
	}
	return nil
}
