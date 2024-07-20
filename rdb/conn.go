package rdb

import (
	"bufio"
	"net"
	"time"
)

type Conn struct {
	conn     net.Conn
	br       *bufio.Reader
	bw       *bufio.Writer
	closed   bool
	idleTime time.Time

	closech chan struct{}
}
