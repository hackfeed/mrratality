package cachedb

import (
	"fmt"
	"os"
	"strings"

	"github.com/tarantool/go-tarantool"
)

var DB *tarantool.Connection

func ConnectDB() {
	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect(os.Getenv("TARANTOOL_URL"), opts)
	if err != nil {
		panic("Failed to connect to Tarantool")
	}

	err = initSpace(conn)
	if err != nil {
		fmt.Println(err)
		panic("Failed to create Tarantool space")
	}

	DB = conn
}

func initSpace(conn *tarantool.Connection) error {
	_, err := conn.Eval("return box.schema.space.create('cache')", []interface{}{})
	if err != nil {
		errString := err.Error()
		if !(strings.Contains(errString, "unsupported") || strings.Contains(errString, "exists")) {
			return err
		}
	}

	_, err = conn.Eval(`return box.space.cache:format({
		{name = 'userfileperiod', type = 'string', is_nullable = false},
		{name = 'mrr', type = 'map', is_nullable = false},
	})`, []interface{}{})
	if err != nil {
		return err
	}

	_, err = conn.Eval("return box.space.cache:create_index('pk', {parts = {'userfileperiod'}})", []interface{}{})
	if err != nil {
		errString := err.Error()
		if !strings.Contains(errString, "exists") {
			return err
		}
	}

	return nil
}
