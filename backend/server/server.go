package server

import (
	"backend/server/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/upload", controllers.SaveFile)
	r.POST("/parse", controllers.ParseFile)
	// r.GET("/load/:id", controllers.LoadData)

	return r
}
