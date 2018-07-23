package git

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

//CloneRepo clones a git repository
func CloneRepo(path, url, branch string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}
	cmdStr := fmt.Sprintf("git clone -b %s %s", branch, url)
	log.Info(cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = path
	return cmd.Run()
}
