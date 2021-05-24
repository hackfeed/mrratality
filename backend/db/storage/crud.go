package storagedb

import "fmt"

func Create(table string, fields []string, data Invoice) error {
	sess := DB.NewSession(nil)
	_, err := sess.InsertInto(table).Columns(fields...).Record(&data).Exec()

	return err
}

func CreateMultiple(table string, fields []string, data []Invoice) error {
	sess := DB.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, invoice := range data {
		_, err := sess.InsertInto(table).Columns(fields...).Record(&invoice).Exec()
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func Read(table string, fields []string) ([]Invoice, error) {
	var data []Invoice

	sess := DB.NewSession(nil)
	_, err := sess.Select(fields...).From(table).LoadStructs(&data)

	return data, err
}

func Delete(table, userID, fileID string) error {
	sess := DB.NewSession(nil)
	_, err := sess.DeleteBySql(
		fmt.Sprintf(
			"ALTER TABLE %v DELETE WHERE user_id = '%v' AND file_id = '%v'", table, userID, fileID,
		)).Exec()

	return err
}
