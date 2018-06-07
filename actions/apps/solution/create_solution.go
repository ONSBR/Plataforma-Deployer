package solution

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//CreateSolution on platform
func CreateSolution(solution *models.Solution) *exceptions.Exception {
	sol := FindSolutionByID(solution.ID)
	if sol != nil {
		return exceptions.NewInvalidArgumentException(fmt.Errorf("solution %s already exists", solution.ID))
	}
	return apicore.PersistOne(solution)
}
