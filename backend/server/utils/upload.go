package utils

import (
	storagedb "backend/db/storage"
	"backend/server/models"
	"fmt"
	"strings"
)

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
