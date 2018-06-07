package env

//GetGitServerKeysPath returns git server keys path
func GetGitServerKeysPath() string {
	return Get("GIT_SERVER_PATH", "/home/philippe/git-server/keys")
}

//GetGitServerReposPath returns git server repos path
func GetGitServerReposPath() string {
	return Get("GIT_SERVER_PATH", "/home/philippe/git-server/repos")
}
