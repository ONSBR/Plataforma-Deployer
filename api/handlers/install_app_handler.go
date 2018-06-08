package handlers

import (
	"github.com/ONSBR/Plataforma-Deployer/actions/apps"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/gin-gonic/gin"
)

//InstallAppHandler handles app instalation
func InstallAppHandler(c *gin.Context) {
	app := models.NewApp()
	c.Bind(app)
	app.SystemID = c.Param("solution")
	if resp, ex := apps.CreateApp(app); ex != nil {
		c.JSON(ex.Status(), ex)
	} else {
		c.JSON(200, resp)
	}
}
