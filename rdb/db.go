package rdb

import "sync"

type RDB struct {
	*ArtTree
	sync.RWMutex
}

func NewRDB() *RDB {
	return &RDB{NewTree(), sync.RWMutex{}}
}

var db = NewRDB()
