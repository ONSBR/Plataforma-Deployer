package git

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/ONSBR/Plataforma-Deployer/env"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCreateAGitRepo(t *testing.T) {
	Convey("should create git repository", t, func() {
		err := CreateGitRepo(fmt.Sprintf("%s/%s", env.GetGitServerReposPath(), "go-repo-git"))
		if err != nil {
			//fmt.Println(err.Message)
			t.Fail()
		}
		exec.Command("bash", fmt.Sprintf("-c rm -rf %s/%s", env.GetGitServerReposPath(), "go-repo-git"))
	})

}
