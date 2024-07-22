package rdb

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"gordb/pkg/logger"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Server struct {
	ln     net.Listener
	laddr  string
	conns  map[*Conn]struct{}
	mu     sync.Mutex
	Closed atomic.Bool
	ctx    context.Context
}

func NewServer(ctx context.Context, addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		ln:    l,
		laddr: l.Addr().String(),
		conns: map[*Conn]struct{}{},
		ctx:   ctx,
	}
	return s, nil
}

func (s *Server) Serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		ctx := context.WithValue(s.ctx, "remoteAddr", conn.RemoteAddr().String())
		go s.serveConn(ctx, conn)
	}
}

func (s *Server) removeConn(c *Conn) {
	s.mu.Lock()
	delete(s.conns, c)
	s.mu.Unlock()
}

func (s *Server) serveConn(ctx context.Context, conn net.Conn) {
	c := &Conn{
		conn:       conn,
		br:         bufio.NewReader(conn),
		bw:         bufio.NewWriter(conn),
		closech:    make(chan struct{}),
		remoteAddr: conn.RemoteAddr().String(),
	}
	s.mu.Lock()
	s.conns[c] = struct{}{}
	s.mu.Unlock()

	readCh := make(chan []string)

	go func() {
		defer close(readCh)
		for {
			select {
			case <-c.closech:
				s.removeConn(c)
				return
			default:
				resp, err := parseRESP(c.br)
				if err != nil {
					// TODO: handle error
					c.err = err
					logger.ErrorWith(ctx, err).Msg("parseRESP error")
					if errors.Is(err, io.EOF) {
						c.Close(ctx)
					}
					break
				}
				readCh <- resp
			}

		}
	}()

	// TODO: handle client requests here
	for cmd := range readCh {
		logger.Info(ctx).Interface("cmd", cmd).Msg("cmd received")

		if handler, ok := routers[cmd[0]]; ok {
			handler(ctx, cmd, &RedisResponse{
				conn: c,
			})
			continue
		}
		c.bw.WriteString(fmt.Sprintf("%cERR unknown command '%s'\r\n", Errors, cmd[0]))
		c.bw.Flush()
	}

}
