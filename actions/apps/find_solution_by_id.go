package apps

import (
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//FindSolutionByID looks for solution on apicore
func FindSolutionByID(id string) (*models.Solution, *exceptions.Exception) {
	list := make([]*models.Solution, 1)
	ex := apicore.FindByID("system", id, &list)
	if len(list) > 0 {
		return list[0], nil
	}
	return new(models.Solution), ex
}
