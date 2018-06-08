package env

import "fmt"

//GetGitServerKeysPath returns git server keys path
func GetGitServerKeysPath() string {
	return Get("GIT_SERVER_PATH", "/home/philippe/git-server/keys")
}

//GetGitServerReposPath returns git server repos path
func GetGitServerReposPath() string {
	return Get("GIT_SERVER_PATH", "/home/philippe/git-server/repos")
}

//GetSSHRemoteURL returns git remote url pattern for ssh protocol
func GetSSHRemoteURL(solution, app string) string {
	user := Get("GIT_SERVER_USER", "git")
	host := Get("GET_SERVER_HOST", "localhost")
	port := Get("GET_SERVER_PORT", "2222")
	return fmt.Sprintf("ssh://%s@%s:%s/git-server/repos/%s/%s", user, host, port, solution, app)
}
