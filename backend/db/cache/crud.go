package cachedb

import "github.com/tarantool/go-tarantool"

func Create(space string, data []interface{}) error {
	_, err := DB.Insert(space, data)

	return err
}

func Read(space, pk string) ([]interface{}, error) {
	resp, err := DB.Select(space, "pk", 0, 1, tarantool.IterEq, []interface{}{pk})
	if err != nil {
		return nil, err
	}

	return resp.Data, err
}
