package rdb

import (
	"bufio"
	"errors"
	"fmt"
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
}

func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		ln:    l,
		laddr: l.Addr().String(),
		conns: map[*Conn]struct{}{},
	}
	return s, nil
}

func (s *Server) Serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.serveConn(conn)
	}
}

func (s *Server) removeConn(c *Conn) {
	s.mu.Lock()
	delete(s.conns, c)
	s.mu.Unlock()
}

func (s *Server) serveConn(conn net.Conn) {
	c := &Conn{
		conn:    conn,
		br:      bufio.NewReader(conn),
		bw:      bufio.NewWriter(conn),
		closech: make(chan struct{}),
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
					fmt.Printf("parseRESP error: %v", err)
					if errors.Is(err, io.EOF) {
						c.Close()
					}
					break
				}
				readCh <- resp
			}

		}
	}()

	// TODO: handle client requests here
	for cmd := range readCh {
		fmt.Printf("cmdLine: %+v \n", cmd)
		c.bw.WriteString("+OK\r\n")
		c.bw.Flush()
	}

}
