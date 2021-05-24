package utils

import (
	storagedb "backend/db/storage"
	userdb "backend/db/user"
	"backend/server/models"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LoadFiles(userID string) ([]map[string]interface{}, error) {
	user, err := userdb.Read(bson.M{"user_id": userID})

	return user.Files, err
}

func UpdateFiles(fname, userID string) error {
	user, err := userdb.Read(bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	var updateObj primitive.D

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	files := append(user.Files, bson.M{"name": fname, "uploaded_at": uploadedAt})
	updateObj = append(updateObj, bson.E{"files", files})

	return userdb.Update(updateObj, bson.M{"user_id": userID})
}

func DeleteFile(userID, fileID string) error {
	err := os.Remove(fmt.Sprintf("static/%v/%v", userID, fileID))
	if err != nil {
		return err
	}

	user, err := userdb.Read(bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	var updateObj primitive.D
	newFiles := []map[string]interface{}{}

	for _, file := range user.Files {
		if file["name"] != fileID {
			newFiles = append(newFiles, file)
		}
	}

	updateObj = append(updateObj, bson.E{"files", newFiles})

	err = userdb.Update(updateObj, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	return storagedb.Delete("mrr.storage", userID, fileID)
}

func UploadFile(userID, fileID string, invoices []*models.Invoice) error {
	var tfInvoices []storagedb.Invoice

	for _, invoice := range invoices {
		tfInvoices = append(tfInvoices, mapModelToEntity(userID, fileID, *invoice))
	}

	return storagedb.CreateMultiple("mrr.storage", storagedb.AllFields, tfInvoices)
}

func mapModelToEntity(userID, fileID string, model models.Invoice) storagedb.Invoice {
	return storagedb.Invoice{
		UserID:      userID,
		FileID:      fileID,
		CustomerId:  model.CustomerId,
		PeriodStart: mapDate(model.PeriodStart),
		PaidPlan:    model.PaidPlan,
		PaidAmount:  model.PaidAmount,
		PeriodEnd:   mapDate(model.PeriodEnd),
	}
}

func mapDate(date string) string {
	dmy := strings.Split(date, ".")

	return fmt.Sprintf("%v-%v-%v", dmy[2], dmy[1], dmy[0])
}
