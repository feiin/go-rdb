package rdb

import (
	"context"
)

func Get(key string) (interface{}, error) {
	return db.Search([]byte(key)), nil
}

func GetString(key string) (string, error) {
	strVal, _ := db.Search([]byte(key)).(string)
	return strVal, nil
}

func Set(key string, value interface{}) bool {
	return db.Insert([]byte(key), value)
}

func SetString(key string, value string) bool {
	return db.Insert([]byte(key), value)
}

func SetStringHandler(ctx context.Context, reqCmdArgs []string, res *RedisResponse) {
	if len(reqCmdArgs) < 3 {
		_ = res.WriteError(NewWrongNumberArgs(reqCmdArgs[0]))
		return
	}

	key := reqCmdArgs[1]
	value := reqCmdArgs[2]
	SetString(key, value)
	_ = res.WriteOK()
}

func GetStringHandler(ctx context.Context, reqCmdArgs []string, res *RedisResponse) {
	if len(reqCmdArgs) != 2 {
		_ = res.WriteError(NewWrongNumberArgs(reqCmdArgs[0]))
		return
	}
	key := reqCmdArgs[1]
	val, err := Get(key)
	if err != nil {
		_ = res.WriteError(err)
		return
	}
	_ = res.WriteBulkStrings(val)
}
