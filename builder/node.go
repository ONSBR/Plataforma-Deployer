package builder

import "os/exec"

type NodeBuilder struct {
}

func (b *NodeBuilder) Build(path string) error {
	cmd := exec.Command("bash", "-c", "npm install")
	cmd.Dir = path
	return cmd.Run()
}
