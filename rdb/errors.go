package rdb

import (
	"errors"
	"fmt"
)

func NewWrongNumberArgs(cmd string) error {
	return errors.New(fmt.Sprintf("ERR wrong number of arguments for '%s' command\r\n", cmd))
}
