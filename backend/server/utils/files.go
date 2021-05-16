package utils

import (
	storagedb "backend/db/storage"
	userdb "backend/db/user"
	"backend/server/models"
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadFiles(userId string) ([]map[string]interface{}, error) {
	var user userdb.User

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	err := userdb.UserCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user.Files, nil
}

func UpdateFiles(fname, userId string) error {
	var user userdb.User

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	err := userdb.UserCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}

	var updateObj primitive.D

	uploadedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	files := append(user.Files, bson.M{"name": fname, "uploaded_at": uploadedAt})
	updateObj = append(updateObj, bson.E{"files", files})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = userdb.UserCol.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		return err
	}

	return nil
}

func Upload(userID string, invoices []*models.Invoice) error {
	sess := storagedb.DB.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, invoice := range invoices {
		err = storagedb.Create(sess, "mrr.storage", storagedb.AllFields, mapModelToEntity(userID, *invoice))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func mapModelToEntity(userID string, model models.Invoice) storagedb.Invoice {
	return storagedb.Invoice{
		UserID:         userID,
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
