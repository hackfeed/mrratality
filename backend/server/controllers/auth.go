package controllers

import (
	userdb "backend/db/user"
	"backend/server/models"
	"backend/server/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	var json models.User

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	count, err := userdb.UserCol.CountDocuments(ctx, bson.M{"email": json.Email})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check email",
		})
		return
	}

	if count > 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "This email is already taken",
		})
		return
	}

	password, err := utils.HashPassword(*json.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	var user userdb.User

	user.Email = json.Email
	user.Password = &password
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, _ := utils.GenerateTokens(*user.Email, user.UserID)
	user.Token = &token
	user.RefreshToken = &refreshToken

	_, err = userdb.UserCol.InsertOne(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	defer cancel()

	expiresAt, err := utils.GetExpirationTime(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User created",
		"idToken":   token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}

func Login(c *gin.Context) {
	var json models.User
	var user userdb.User

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = userdb.UserCol.FindOne(ctx, bson.M{"email": json.Email}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "User with given email doesn't exist",
		})
		return
	}

	err = utils.VerifyPassword(*json.Password, *user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Password is incorrect",
		})
		return
	}

	token, refreshToken, err := utils.GenerateTokens(*user.Email, user.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate tokens",
		})
		return
	}

	utils.UpdateTokens(token, refreshToken, user.UserID)

	expiresAt, err := utils.GetExpirationTime(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get token expiration time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"idToken":   token,
		"localId":   user.UserID,
		"expiresAt": expiresAt,
	})
}
