package handlers

import (
	"github.com/ONSBR/Plataforma-Deployer/actions/apps/solution"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/gin-gonic/gin"
)

//CreateSolutionHandler handle create solution service
func CreateSolutionHandler(c *gin.Context) {
	sol := models.NewSolution()
	err := c.BindJSON(sol)
	if err != nil {
		ex := exceptions.NewInvalidArgumentException(err)
		c.JSON(ex.Status(), ex)
		return
	}
	ex := solution.CreateSolution(sol)
	if ex != nil {
		c.JSON(ex.Status(), ex)
		return
	}
	c.JSON(201, gin.H{"message": "solution created"})
}
