package handlers

import (
	"github.com/ONSBR/Plataforma-Deployer/actions/apps"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/gin-gonic/gin"
)

//CreateSolutionHandler handle create solution service
func CreateSolutionHandler(c *gin.Context) {
	sol := models.NewSolution()
	err := c.BindJSON(sol)
	if err != nil {
		c.JSON(500, err)
		return
	}
	ex := apps.CreateSolution(sol)
	if ex != nil {
		c.JSON(400, ex)
		return
	}
	c.JSON(201, gin.H{"message": "solution created"})
}
