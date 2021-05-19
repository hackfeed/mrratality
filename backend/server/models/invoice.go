package models

type Invoice struct {
	InvoiceCreated string  `csv:"invoice_created"`
	InvoiceId      uint32  `csv:"invoice_id"`
	CustomerId     uint32  `csv:"customer_id"`
	PaidAmount     float32 `csv:"paid_amount"`
	PaidCurrency   string  `csv:"paid_currency"`
	PeriodStart    string  `csv:"period_start"`
	PeriodEnd      string  `csv:"period_end"`
	PaidUser       uint32  `csv:"paid_user"`
	PaidPlan       string  `csv:"paid_plan"`
}
