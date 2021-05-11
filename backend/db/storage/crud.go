package storagedb

import (
	"github.com/mailru/dbr"
)

func Create(sess *dbr.Session, table string, fields []string, data Invoice) error {
	_, err := sess.InsertInto(table).Columns(fields...).Record(&data).Exec()
	if err != nil {
		return err
	}
	return nil
}

func Read(sess *dbr.Session, table string, fields []string) ([]Invoice, error) {
	var data []Invoice
	_, err := sess.Select(fields...).From(table).LoadStructs(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
