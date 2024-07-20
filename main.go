package main

import (
	"fmt"
	"gordb/rdb"
)

func main() {
	addr := ":6378"
	s, err := rdb.NewServer(addr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server running on %s", addr)
	s.Serve()
	fmt.Println("Server stopped")
}
