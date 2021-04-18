package server

import "github.com/gin-gonic/gin"

func SetupServer() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", PingEndpoint)

	return r
}

func PingEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
