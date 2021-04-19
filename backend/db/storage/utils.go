package storagedb

import (
	"log"

	"github.com/mailru/dbr"
)

func InitDB(conn *dbr.Connection) error {
	sess := conn.NewSession(nil)
	_, err := sess.Exec(`
	CREATE DATABASE IF NOT EXISTS mrr
	`)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func InitTables(conn *dbr.Connection) error {
	sess := conn.NewSession(nil)
	_, err := sess.Exec(`
	CREATE TABLE IF NOT EXISTS mrr.storage(
		invoice_created Date,
		invoice_id UInt32,
		customer_id UInt32,
		paid_amount UInt32,
		paid_currency FixedString(3),
		period_start Date,
		period_end Date,
		paid_users Decimal32(2),
		paid_plan String
	) 
	Engine=Memory
	`)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Init(conn *dbr.Connection) error {
	err := InitDB(conn)
	if err != nil {
		log.Println(err)
		return err
	}
	err = InitTables(conn)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
