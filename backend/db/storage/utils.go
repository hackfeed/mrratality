package storagedb

import (
	"os"

	"github.com/mailru/dbr"
	_ "github.com/mailru/go-clickhouse"
)

var DB *dbr.Connection

func ConnectDB() {
	conn, err := dbr.Open("clickhouse", os.Getenv("CLICKHOUSE_URL"), nil)
	if err != nil {
		panic("Failed to connect to Clickhouse")
	}

	err = initDB(conn)
	if err != nil {
		panic("Failed to create Clickhouse database")
	}

	err = initTable(conn)
	if err != nil {
		panic("Failed to create Clickhouse table")
	}

	DB = conn
}

func initDB(conn *dbr.Connection) error {
	sess := conn.NewSession(nil)
	_, err := sess.Exec(`
	CREATE DATABASE IF NOT EXISTS mrr
	`)
	if err != nil {
		return err
	}

	return nil
}

func initTable(conn *dbr.Connection) error {
	sess := conn.NewSession(nil)
	_, err := sess.Exec(`
	CREATE TABLE IF NOT EXISTS mrr.storage(
		user_id String,
		invoice_created Date,
		invoice_id UInt32,
		customer_id UInt32,
		paid_amount Float32,
		paid_currency FixedString(3),
		period_start Date,
		period_end Date,
		paid_user UInt32,
		paid_plan String
	) 
	Engine=Memory
	`)
	if err != nil {
		return err
	}

	return nil
}
