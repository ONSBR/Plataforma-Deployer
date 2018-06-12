package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

//CloneRepo clones a git repository
func CloneRepo(path, url, branch string) *exceptions.Exception {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && git clone -b %s %s", path, branch, url))
	cmd.Dir = path
	return exceptions.NewComponentException(cmd.Run())
}
