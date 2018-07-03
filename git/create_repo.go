package git

import (
	"fmt"
	"os"
	"os/exec"
)

//CreateGitRepo creates a new bare git repository on git-server
func CreateGitRepo(path string) error {
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		return fmt.Errorf("repository %s already exist", path)
	}
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}
	cmd := exec.Command("bash", "-c", "git init --bare --shared=true")
	cmd.Dir = path
	return cmd.Run()
}
