package utils

import (
	storagedb "backend/db/storage"
	userdb "backend/db/user"
	"backend/server/models"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LoadFiles(userId string) ([]map[string]interface{}, error) {
	user, err := userdb.Read(bson.M{"user_id": userId})
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return user.Files, nil
}

func UpdateFiles(fname, userId string) error {
	user, err := userdb.Read(bson.M{"user_id": userId})
	if err != nil {
		return err
	}

	var updateObj primitive.D

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	files := append(user.Files, bson.M{"name": fname, "uploaded_at": uploadedAt})
	updateObj = append(updateObj, bson.E{"files", files})

	return userdb.Update(updateObj, bson.M{"user_id": userId})
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
		UserID:         userID,
		FileID:         fileID,
		InvoiceCreated: mapDate(model.InvoiceCreated),
		InvoiceId:      model.InvoiceId,
		CustomerId:     model.CustomerId,
		PaidAmount:     model.PaidAmount,
		PaidCurrency:   model.PaidCurrency,
		PeriodStart:    mapDate(model.PeriodStart),
		PeriodEnd:      mapDate(model.PeriodEnd),
		PaidUser:       model.PaidUser,
		PaidPlan:       model.PaidPlan,
	}
}

func mapDate(date string) string {
	dmy := strings.Split(date, ".")

	return fmt.Sprintf("%v-%v-%v", dmy[2], dmy[1], dmy[0])
}
