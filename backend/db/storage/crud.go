package storagedb

import (
	"log"

	"github.com/mailru/dbr"
)

func Create(conn *dbr.Connection, table string, fields []string, data StorageEntity) error {
	sess := conn.NewSession(nil)
	_, err := sess.InsertInto(table).Columns(fields...).Record(&data).Exec()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Read(conn *dbr.Connection, table string, fields []string) ([]StorageEntity, error) {
	sess := conn.NewSession(nil)
	var data []StorageEntity
	_, err := sess.Select(fields...).From(table).LoadStructs(&data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil
}
