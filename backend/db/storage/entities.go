package storagedb

import (
	"time"
)

type StorageEntity struct {
	InvoiceCreated time.Time `db:"invoice_created"`
	InvoiceId      uint32    `db:"invoice_id"`
	CustomerId     uint32    `db:"customer_id"`
	PaidAmount     uint32    `db:"paid_amount"`
	PaidCurrency   string    `db:"paid_currency"`
	PeriodStart    time.Time `db:"period_start"`
	PeriodEnd      time.Time `db:"period_end"`
	PaidUsers      float32   `db:"paid_users"`
	PaidPlan       string    `db:"paid_plan"`
}
