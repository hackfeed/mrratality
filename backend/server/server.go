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

	r.Use(cors.Default())
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.Use(middlewares.Auth())
	r.POST("/upload", controllers.SaveFile)
	r.POST("/parse", controllers.ParseFile)
	// r.GET("/load/:id", controllers.LoadData)

	return r
}
