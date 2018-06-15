package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	log "github.com/sirupsen/logrus"
)

//CloneRepo clones a git repository
func CloneRepo(path, url, branch string) *exceptions.Exception {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	cmdStr := fmt.Sprintf("git clone -b %s %s", branch, url)
	log.Info(cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = path
	return exceptions.NewComponentException(cmd.Run())
}
