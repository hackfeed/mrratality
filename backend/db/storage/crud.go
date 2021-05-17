package storagedb

func Create(table string, fields []string, data Invoice) error {
	sess := DB.NewSession(nil)
	_, err := sess.InsertInto(table).Columns(fields...).Record(&data).Exec()
	if err != nil {
		return err
	}

	return nil
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
	if err != nil {
		return data, err
	}

	return data, nil
}
