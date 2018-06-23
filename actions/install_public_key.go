package actions

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

//InstallPublicKey install user's public key on git-server
func InstallPublicKey(content []byte, solution, keyName string) (*models.PublicKeyInfo, *exceptions.Exception) {
	fd, err := os.Create(fmt.Sprintf("%s/%s", env.GetGitServerKeysPath(), keyName))
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	_, err = fd.Write(content)
	if err != nil {
		return nil, exceptions.NewComponentException(err)
	}
	fd.Close()
	info := new(models.PublicKeyInfo)
	info.Name = keyName
	info.Solution = solution
	info.Ok()

	return info, nil
}
