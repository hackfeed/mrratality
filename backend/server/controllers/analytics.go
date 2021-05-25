package controllers

import (
	"backend/server/models"
	"backend/server/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(c *gin.Context) {
	userID, ie := c.Get("userID")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	var req models.Period

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	data, err := utils.GetAnalytics(userID.(string), req.Filename, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Analytics is loaded",
		"data":    data,
	})
}
