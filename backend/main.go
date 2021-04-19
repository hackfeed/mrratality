package main

import (
	storagedb "backend/db/storage"
	server "backend/server"
	"fmt"
	"log"
	"time"

	"github.com/mailru/dbr"
	_ "github.com/mailru/go-clickhouse"
)

func main() {
	conn, err := dbr.Open("clickhouse", "http://localhost:8123/default", nil)
	if err != nil {
		log.Fatal(err)
	}
	err = storagedb.Init(conn)
	if err != nil {
		log.Fatal(err)
	}
	data := storagedb.StorageEntity{
		InvoiceCreated: time.Now(),
		InvoiceId:      2,
		CustomerId:     1,
		PaidAmount:     1,
		PaidCurrency:   "EUR",
		PeriodStart:    time.Now(),
		PeriodEnd:      time.Now(),
		PaidUsers:      2.0,
		PaidPlan:       "cool",
	}
	fields := []string{
		"invoice_created",
		"invoice_id",
		"customer_id",
		"paid_amount",
		"paid_currency",
		"period_start",
		"period_end",
		"paid_users",
		"paid_plan",
	}
	err = storagedb.Create(conn, "mrr.storage", fields, data)
	if err != nil {
		log.Fatal(err)
	}
	rd, err := storagedb.Read(conn, "mrr.storage", fields)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rd)
	server.SetupServer().Run()
}
