package api

import (
	"github.com/gin-gonic/gin"
)

//Run starts API listen on default port
func Run() {
	router := gin.Default()

	router.Run(":6970")
}
