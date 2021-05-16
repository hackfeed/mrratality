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
			"files":   map[string]interface{}{},
		})
		return
	}

	files, err := utils.LoadFiles(userId.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user files",
			"files":   map[string]interface{}{},
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
			"message":  "No file is received",
			"filename": "",
		})
		return
	}

	fext := filepath.Ext(file.Filename)
	if fext != ".csv" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":  "Wrong file format. Please provide CSV file",
			"filename": "",
		})
		return
	}

	userId, ie := c.Get("userId")
	if !ie {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":  "Unable to determine logged in user",
			"filename": "",
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
			"message":  "Unable to save the file",
			"filename": "",
		})
		return
	}

	err = utils.UpdateFiles(fname, userId.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":  "Unable to update user in db",
			"filename": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded",
		"filename": fname,
	})
}

func ParseFile(c *gin.Context) {
	var req models.ParseRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	file, err := os.Open(fmt.Sprintf("static/%v.csv", req.UUID))
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

	err = utils.Upload(req.UUID, invoices)
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
