package handlers

import "github.com/gin-gonic/gin"

//DeployAppHandler starts deploy process
func DeployAppHandler(c *gin.Context) {
	c.AbortWithStatus(500)
}
