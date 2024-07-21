package rdb

import (
	"bufio"
	"context"
	"gordb/pkg/logger"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	conn       net.Conn
	br         *bufio.Reader
	bw         *bufio.Writer
	closed     atomic.Bool
	idleTime   time.Time
	err        error
	remoteAddr string

	closech chan struct{}
}

func (c *Conn) Err() error {
	return c.err
}

func (c *Conn) Close(ctx context.Context) error {
	if c.closed.Load() {
		return nil
	}

	c.cleanup(ctx)
	return nil
}

func (c *Conn) cleanup(ctx context.Context) {
	if c.closed.Swap(true) {
		return
	}
	close(c.closech)
	if c.conn == nil {
		return
	}
	err := c.conn.Close()
	if err != nil {
		logger.ErrorWith(ctx, err).Msg("conn cleanup error")
	}
}
