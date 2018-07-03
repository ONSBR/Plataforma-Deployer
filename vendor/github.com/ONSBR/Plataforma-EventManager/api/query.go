package api

import (
	"github.com/ONSBR/Plataforma-EventManager/eventstore"
	"github.com/gin-gonic/gin"
)

func registerQueryEndpoints(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/count", func(c *gin.Context) {
		field := c.Query("field")
		value := c.Query("value")
		last := c.Query("last")
		c.JSON(200, gin.H{
			"total": eventstore.Count(field, value, last),
		})
	})

	r.GET("/events", func(c *gin.Context) {
		field := c.Query("field")
		value := c.Query("value")
		last := c.Query("last")
		c.JSON(200, gin.H{
			"result": eventstore.Query(field, value, last),
		})
	})
}
