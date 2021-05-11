package server

import (
	storagedb "backend/db/storage"
	"backend/server/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	storagedb.ConnectDB()

	r.Use(cors.Default())

	r.POST("/upload", controllers.SaveFile)
	r.POST("/parse", controllers.ParseFile)
	// r.GET("/load/:id", controllers.LoadData)

	return r
}
