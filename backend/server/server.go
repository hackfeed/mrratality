package server

import (
	storagedb "backend/db/storage"
	userdb "backend/db/user"
	"backend/server/controllers"
	"backend/server/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	storagedb.ConnectDB()
	userdb.ConnectDB()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"token"}
	r.Use(cors.New(config))

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Use(middlewares.Auth())

	files := r.Group("/files")
	{
		files.POST("/upload", controllers.SaveFile)
		files.POST("/parse", controllers.ParseFile)
		files.GET("/load", controllers.LoadFiles)
		// r.GET("/load/:id", controllers.LoadData)
	}

	return r
}
