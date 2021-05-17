package storagedb

var AllFields = []string{
	"user_id",
	"file_id",
	"invoice_created",
	"invoice_id",
	"customer_id",
	"paid_amount",
	"paid_currency",
	"period_start",
	"period_end",
	"paid_user",
	"paid_plan",
}

type Invoice struct {
	UserID         string  `db:"user_id"`
	FileID         string  `db:"file_id"`
	InvoiceCreated string  `db:"invoice_created"`
	InvoiceId      uint32  `db:"invoice_id"`
	CustomerId     uint32  `db:"customer_id"`
	PaidAmount     float32 `db:"paid_amount"`
	PaidCurrency   string  `db:"paid_currency"`
	PeriodStart    string  `db:"period_start"`
	PeriodEnd      string  `db:"period_end"`
	PaidUser       uint32  `db:"paid_user"`
	PaidPlan       string  `db:"paid_plan"`
}
