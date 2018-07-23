package apps

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//FindAppByID on apicore
func FindAppByID(id string) (*models.App, error) {
	list := make([]*models.App, 1)
	if ex := apicore.FindByID(models.NewApp().Metadata.Type, id, &list); ex != nil {
		return nil, ex
	} else if len(list) > 0 {
		return list[0], nil
	}
	return nil, fmt.Errorf("App with id %s not found on apicore", id)
}
