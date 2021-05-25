package storagedb

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

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

func CreateMPPTable(conn *dbr.Connection, mpp MPP) (string, error) {
	tmpl :=
		`
	CREATE TABLE IF NOT EXISTS mrr."{{ .UserFileID }}-{{ .PeriodStart }}-{{ .PeriodEnd }}"(
		userfile_id String,
		{{ $first := true }}{{ range .Dates }}{{if $first}}{{$first = false}}{{else}},{{end}}
		"{{ . }}" Float32{{ end }}
	)
	ENGINE = SummingMergeTree()
	ORDER BY userfile_id
	`

	t, err := template.New("crtbl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, mpp)
	if err != nil {
		return "", err
	}

	query := buf.String()
	fmt.Println(query)
	sess := conn.NewSession(nil)
	_, err = sess.Exec(query)
	if err != nil {
		return "", err
	}

	return "", nil
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
	ENGINE = Memory
	`)
	if err != nil {
		return err
	}

	return nil
}
