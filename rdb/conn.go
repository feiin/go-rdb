package rdb

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	conn     net.Conn
	br       *bufio.Reader
	bw       *bufio.Writer
	closed   atomic.Bool
	idleTime time.Time
	err      error

	closech chan struct{}
}

func (c *Conn) Err() error {
	return c.err
}

func (c *Conn) Close() error {
	if c.closed.Load() {
		return nil
	}

	c.cleanup()
	return nil
}

func (c *Conn) cleanup() {
	if c.closed.Swap(true) {
		return
	}
	close(c.closech)
	if c.conn == nil {
		return
	}
	err := c.conn.Close()
	if err != nil {
		fmt.Printf("conn close error: %v", err)
	}
}
