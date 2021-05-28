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

func LoadFiles(c *gin.Context) {
	userID, ie := c.Get("userID")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	files, err := utils.LoadFiles(userID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded",
		"files":   files,
	})
}

func DeleteFile(c *gin.Context) {
	userID, ie := c.Get("userID")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	var req models.File

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = utils.DeleteFile(userID.(string), req.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted",
	})
}

func SaveFile(c *gin.Context) {
	userID, ie := c.Get("userID")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	fext := filepath.Ext(file.Filename)
	if fext != ".csv" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Wrong file format. Please provide CSV file",
		})
		return
	}

	dir := fmt.Sprintf("static/%v", userID)
	fname := fmt.Sprintf("%v%v", uuid.New(), fext)
	fpath := fmt.Sprintf("%v/%v", dir, fname)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}

	err = c.SaveUploadedFile(file, fpath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	err = utils.UpdateFiles(fname, userID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update user in db",
		})
		return
	}

	csvFile, err := os.Open(fpath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to find data with given id",
		})
		return
	}
	defer csvFile.Close()

	invoices := []*models.Invoice{}

	err = gocsv.UnmarshalFile(csvFile, &invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse given CSV file",
		})
		return
	}

	err = utils.UploadFile(userID.(string), fname, invoices)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to upload data to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded",
		"filename": fname,
	})
}
