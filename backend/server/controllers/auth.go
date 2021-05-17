package controllers

import (
	userdb "backend/db/user"
	"backend/server/models"
	"backend/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	user, err := userdb.SignUp(*req.Email, *req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	expiresAt, err := utils.GetExpirationTime(*user.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User created",
		"idToken":   *user.Token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}

func Login(c *gin.Context) {
	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	user, err := userdb.Login(*req.Email, *req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	expiresAt, err := utils.GetExpirationTime(*user.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"idToken":   *user.Token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}
