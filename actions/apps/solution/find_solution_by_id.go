package solution

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//FindSolutionByID looks for solution on apicore
func FindSolutionByID(id string) *models.Solution {
	list := make([]*models.Solution, 1)
	apicore.FindByID("system", id, &list)
	if len(list) > 0 {
		return list[0]
	}
	return nil
}
