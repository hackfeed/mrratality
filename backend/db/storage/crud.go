package storagedb

import (
	"fmt"
	"strings"

	"github.com/mailru/dbr"
)

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

func CreateDynamic(table string, data []interface{}) error {
	vals := ""
	for _, val := range data {
		vals += fmt.Sprintf("%v,", val)
	}
	vals = strings.TrimSuffix(vals, ",")

	sess := DB.NewSession(nil)
	_, err := sess.InsertBySql(fmt.Sprintf(
		"INSERT INTO %v (*) VALUES (%v)", table, vals,
	)).Exec()

	return err
}

func Read(table string, fields []string) ([]Invoice, error) {
	var data []Invoice

	sess := DB.NewSession(nil)
	_, err := sess.Select(fields...).From(table).LoadStructs(&data)

	return data, err
}

func ReadByPeriod(table string, fields []string, userID, fileID, periodStart, periodEnd string) ([]Invoice, error) {
	var data []Invoice

	sess := DB.NewSession(nil)
	_, err := sess.Select(fields...).From(table).Where(
		dbr.And(
			dbr.Lte("period_start", periodEnd),
			dbr.Gte("period_end", periodStart),
			dbr.Eq("user_id", userID),
			dbr.Eq("file_id", fileID),
		)).LoadStructs(&data)

	return data, err
}

func ReadDynamic(table string) ([][]interface{}, error) {
	rows, err := DB.Query(fmt.Sprintf("SELECT COLUMNS('\\.\\d{4}') FROM %v", table))
	if err != nil {
		return nil, err
	}

	res := [][]interface{}{}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(cols)
	vals := make([]interface{}, count)
	valsPtrs := make([]interface{}, count)

	for i := range cols {
		valsPtrs[i] = &vals[i]
	}

	for rows.Next() {
		err = rows.Scan(valsPtrs...)
		if err != nil {
			return nil, err
		}
		row := []interface{}{}
		row = append(row, vals...)
		res = append(res, row)
	}

	return res, nil
}

func Delete(table, userID, fileID string) error {
	sess := DB.NewSession(nil)
	_, err := sess.DeleteBySql(
		fmt.Sprintf(
			"ALTER TABLE %v DELETE WHERE user_id = '%v' AND file_id = '%v'", table, userID, fileID,
		)).Exec()

	return err
}
