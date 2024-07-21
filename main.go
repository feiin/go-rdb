package main

import (
	"context"
	"fmt"
	"gordb/pkg/logger"
	"gordb/rdb"
)

func main() {
	addr := ":6378"
	ctx := context.WithValue(context.Background(), "name", "gordb")
	s, err := rdb.NewServer(ctx, addr)
	if err != nil {
		fmt.Println(err)
	}
	logger.Info(ctx).Str("addr", addr).Msg("Server is running")
	s.Serve()
	logger.Info(ctx).Str("addr", addr).Msg("Server stopped")
}
