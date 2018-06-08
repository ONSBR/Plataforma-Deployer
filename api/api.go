package api

import (
	"github.com/ONSBR/Plataforma-Deployer/api/handlers"
	"github.com/gin-gonic/gin"
)

//Run starts API listen on default port
func Run() {
	router := gin.Default()
	group := router.Group("api/v1.0.0")

	group.POST("/publickey/:solution/:filename", handlers.UploadPublicKey)
	group.POST("/solution", handlers.CreateSolutionHandler)
	group.POST("/solution/:solution/create/app", handlers.InstallAppHandler)
	router.Run(":6970")
}
