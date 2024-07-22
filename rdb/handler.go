package rdb

import "context"

type handler func(ctx context.Context, reqCmdArgs []string, res *RedisResponse)

var routers = map[string]handler{}

func init() {
	routers["set"] = SetStringHandler
	routers["get"] = GetStringHandler
}
