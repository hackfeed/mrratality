package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/server/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "No Authorization header provided",
			})
			return
		}

		claims, err := utils.ValidateToken(clientToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Token validation failed",
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
