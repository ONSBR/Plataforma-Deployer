package handlers

import (
	"github.com/ONSBR/Plataforma-Deployer/actions/apps"
	"github.com/gin-gonic/gin"
)

//DeployAppHandler starts deploy process
func DeployAppHandler(c *gin.Context) {
	processID := c.Param("processID")
	if ex := apps.CreateDeploy(processID); ex != nil {
		c.JSON(ex.Status(), ex)
	} else {
		c.Status(202)
	}
}
