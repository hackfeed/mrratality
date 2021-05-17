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
	userId, ie := c.Get("userId")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	files, err := utils.LoadFiles(userId.(string))
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

func SaveFile(c *gin.Context) {
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

	userId, ie := c.Get("userId")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to determine logged in user",
		})
		return
	}

	dir := fmt.Sprintf("static/%v", userId)
	fname := fmt.Sprintf("%v%v", uuid.New(), fext)
	fpath := fmt.Sprintf("%v/%v", dir, fname)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}

	err = c.SaveUploadedFile(file, fpath)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	err = utils.UpdateFiles(fname, userId.(string))
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

	err = utils.UploadFile(userId.(string), fname, invoices)
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
