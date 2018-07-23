package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/PMoneda/whaler"
)

//InstallPublicKey install user's public key on git-server
func InstallPublicKey(content []byte, solution, keyName string) (*models.PublicKeyInfo, error) {
	fd, err := os.Create(fmt.Sprintf("%s/%s", env.GetGitServerKeysPath(), keyName))
	if err != nil {
		return nil, err
	}
	_, err = fd.Write(content)
	if err != nil {
		return nil, err
	}
	fd.Close()
	info := new(models.PublicKeyInfo)
	info.Name = keyName
	info.Solution = solution
	timeout := 60 * time.Second
	if err := whaler.RestartContainer("/git-server", &timeout); err != nil {
		return nil, err
	}
	info.Ok()
	return info, nil
}
