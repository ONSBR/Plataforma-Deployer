package models

import (
	"fmt"
)

//DeployContext is the entity to manage all deploy steps
type DeployContext struct {
	Info     *Deploy
	Solution Solution
	RootPath string
	Version  string
}

//GetImageTag returns docker image name pattern
func (d DeployContext) GetImageTag() string {
	return fmt.Sprintf("%s/%s:%s", d.Solution.Name, d.Info.Name, d.Version)
}
