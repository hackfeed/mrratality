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
		file_id String,
		customer_id UInt32,
		period_start Date,
		paid_plan String,
		paid_amount Float32,
		period_end Date
	) 
	Engine=Memory
	`)
	if err != nil {
		return err
	}

	return nil
}
