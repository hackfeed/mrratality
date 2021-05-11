package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"

	"backend/server/models"
	"backend/server/utils"
)

func SaveFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "No file is received",
			"isUploaded": false,
		})
		return
	}

	fext := filepath.Ext(file.Filename)
	if fext != ".csv" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Wrong file format. Please provide CSV file",
			"isUploaded": false,
		})
		return
	}

	err = c.SaveUploadedFile(file, fmt.Sprintf("static/%v%v", uuid.New().String(), fext))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to save the file",
			"isUploaded": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "File uploaded",
		"isUploaded": true,
	})
}

func ParseFile(c *gin.Context) {
	var json models.ParseRequest

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	file, err := os.Open(fmt.Sprintf("static/%v.csv", json.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to find data with given id",
		})
		return
	}
	defer file.Close()

	invoices := []*models.Invoice{}

	err = gocsv.UnmarshalFile(file, &invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse given CSV file",
		})
		return
	}

	err = utils.Upload(json.UUID, invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to upload data to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File parsed",
	})
}