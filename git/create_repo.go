package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

//CreateGitRepo creates a new bare git repository on git-server
func CreateGitRepo(path string) *exceptions.Exception {
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		return exceptions.NewInvalidArgumentException(fmt.Errorf("repository %s already exist", path))
	}
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	cmd := exec.Command("bash", "-c", "git init --bare --shared=true")
	cmd.Dir = path
	return exceptions.NewComponentException(cmd.Run())
}
